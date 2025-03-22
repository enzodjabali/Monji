package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Global variables
var (
	mongoClient *mongo.Client
	db          *sql.DB                    // SQLite database connection
	jwtSecret   = []byte("supersecretkey") // In production, set this via an environment variable.
)

// User represents our user model.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Company   string `json:"company,omitempty"`
	Password  string `json:"password"` // Accept JSON input; we'll clear it before responses.
	Role      string `json:"role"`     // "user", "admin", or "superadmin"
}

// Environment represents a MongoDB environment.
type Environment struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ConnectionString string `json:"connection_string"`
	CreatedBy        int    `json:"created_by"`
}

// initSQLite opens/creates the SQLite DB at /data/sqlite/users.db,
// creates necessary tables, and inserts a default admin user if none exist.
func initSQLite() {
	var err error
	// Open SQLite DB from the persistent volume.
	db, err = sql.Open("sqlite3", "/data/sqlite/users.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// Create the users table if it doesn't exist.
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		company TEXT,
		password TEXT NOT NULL,
		role TEXT NOT NULL
	);
	`
	_, err = db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create the environments table if it doesn't exist.
	createEnvsTableSQL := `
	CREATE TABLE IF NOT EXISTS environments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		connection_string TEXT NOT NULL,
		created_by INTEGER NOT NULL
	);
	`
	_, err = db.Exec(createEnvsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create environments table: %v", err)
	}

	// Check if any user exists.
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check users table: %v", err)
	}
	if count == 0 {
		// Create default admin user.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}
		_, err = db.Exec(
			`INSERT INTO users(first_name, last_name, email, company, password, role)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			"Admin", "User", "admin@example.com", "", string(hashedPassword), "admin")
		if err != nil {
			log.Fatalf("Failed to insert default admin user: %v", err)
		}
		log.Println("Default admin user created: email=admin@example.com, password=admin")
	}
}

func main() {
	// Initialize SQLite (for user auth & environments).
	initSQLite()

	// Connect to MongoDB (the main MongoDB connection for your application).
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:strongpassword@mongodb:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create a Gin router.
	router := gin.Default()

	// Auth endpoints.
	// Note: The /register endpoint is removed. Users are now pre-populated.
	router.POST("/login", loginUser)

	// Environment endpoints:
	// Only admin or superadmin can create new MongoDB environments.
	router.POST("/environments", AuthMiddleware(), AdminMiddleware(), createEnvironment)
	// Any authenticated user can list environments.
	router.GET("/environments", AuthMiddleware(), listEnvironments)

	// Mongo queries using an environment.
	router.GET("/environments/:id/databases", AuthMiddleware(), getDatabasesForEnv)
	router.GET("/environments/:id/databases/:dbName/collections", AuthMiddleware(), getCollectionsForEnv)

	// Start the API server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// loginUser logs in a user and returns a JWT token.
func loginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user from SQLite.
	var user User
	row := db.QueryRow(`SELECT id, first_name, last_name, email, company, password, role
	                    FROM users WHERE email = ?`, credentials.Email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Password, &user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare provided password with stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create a JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(), // expires in 72 hours
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// -------------------------------------------------------------------
// Middleware: Auth & Role Checks
// -------------------------------------------------------------------

// AuthMiddleware verifies the JWT token and loads the user into the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}
			var user User
			row := db.QueryRow(`SELECT id, first_name, last_name, email, company, password, role
			                    FROM users WHERE id = ?`, int(userID))
			err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Password, &user.Role)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			user.Password = ""
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
	}
}

// AdminMiddleware ensures that the user has the role "admin" or "superadmin".
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		user, ok := userVal.(User)
		if !ok || (user.Role != "admin" && user.Role != "superadmin") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - Admins Only"})
			return
		}
		c.Next()
	}
}

// -------------------------------------------------------------------
// Environment Endpoints
// -------------------------------------------------------------------

// createEnvironment allows an admin/superadmin to create a new MongoDB environment.
func createEnvironment(c *gin.Context) {
	var req struct {
		Name             string `json:"name"`
		ConnectionString string `json:"connection_string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" || req.ConnectionString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing environment name or connection string"})
		return
	}

	// Get the current user from context.
	userVal, _ := c.Get("user")
	user := userVal.(User)

	// Insert the new environment into the DB.
	stmt, err := db.Prepare(`INSERT INTO environments (name, connection_string, created_by) VALUES (?,?,?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	res, err := stmt.Exec(req.Name, req.ConnectionString, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last insert id"})
		return
	}

	env := Environment{
		ID:               int(id),
		Name:             req.Name,
		ConnectionString: req.ConnectionString,
		CreatedBy:        user.ID,
	}
	c.JSON(http.StatusOK, gin.H{"environment": env})
}

// listEnvironments returns all stored MongoDB environments.
func listEnvironments(c *gin.Context) {
	rows, err := db.Query(`SELECT id, name, connection_string, created_by FROM environments`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var envs []Environment
	for rows.Next() {
		var e Environment
		if err := rows.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		envs = append(envs, e)
	}
	c.JSON(http.StatusOK, gin.H{"environments": envs})
}

// -------------------------------------------------------------------
// Mongo Queries by Environment
// -------------------------------------------------------------------

// getDatabasesForEnv connects to the specified environment's MongoDB and lists its databases.
func getDatabasesForEnv(c *gin.Context) {
	envIDStr := c.Param("id")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Retrieve the environment from the DB.
	var e Environment
	row := db.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	err = row.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Connect to the environment's MongoDB.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(e.ConnectionString))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	result, err := client.ListDatabases(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// getCollectionsForEnv lists all collections and their stats for a specific database
// in the given environment.
func getCollectionsForEnv(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Retrieve the environment.
	var e Environment
	row := db.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	err = row.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Connect to that MongoDB instance.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(e.ConnectionString))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	mongoDB := client.Database(dbName)
	collections, err := mongoDB.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var output []gin.H
	for _, coll := range collections {
		stats := bson.M{}
		err := mongoDB.RunCommand(ctx, bson.D{{Key: "collStats", Value: coll}}).Decode(&stats)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get stats for %s: %v", coll, err)})
			return
		}

		indexes, err := mongoDB.Collection(coll).Indexes().List(ctx)
		var indexList []bson.M
		if err == nil {
			for indexes.Next(ctx) {
				var idx bson.M
				if err := indexes.Decode(&idx); err == nil {
					indexList = append(indexList, idx)
				}
			}
		}

		output = append(output, gin.H{
			"name":           coll,
			"count":          stats["count"],
			"size":           stats["size"],
			"storageSize":    stats["storageSize"],
			"totalIndexSize": stats["totalIndexSize"],
			"indexes":        indexList,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"database":    dbName,
		"collections": output,
	})
}
