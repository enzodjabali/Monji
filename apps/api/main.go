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

// -------------------------------------------------------------------
// Global variables & data structures
// -------------------------------------------------------------------

var (
	db        *sql.DB                    // SQLite database connection
	jwtSecret = []byte("supersecretkey") // Ideally set via an ENV var
)

// User represents a row in the "users" table.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Company   string `json:"company,omitempty"`
	Password  string `json:"password"` // allow JSON input on register; cleared before returning
	Role      string `json:"role"`     // "user", "admin", or "superadmin"
}

// Environment represents a row in the "environments" table.
type Environment struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ConnectionString string `json:"connection_string"`
	CreatedBy        int    `json:"created_by"`
}

// -------------------------------------------------------------------
// Initialization
// -------------------------------------------------------------------

func initSQLite() {
	var err error
	// Open or create local SQLite DB file: "users.db"
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// 1) Create the users table if it doesn't exist.
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
	if _, err = db.Exec(createUsersTableSQL); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// 2) Create the environments table if it doesn't exist.
	createEnvsTableSQL := `
	CREATE TABLE IF NOT EXISTS environments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		connection_string TEXT NOT NULL,
		created_by INTEGER NOT NULL
	);
	`
	if _, err = db.Exec(createEnvsTableSQL); err != nil {
		log.Fatalf("Failed to create environments table: %v", err)
	}
}

// -------------------------------------------------------------------
// Main Entry Point
// -------------------------------------------------------------------

func main() {
	// Initialize SQLite (for user auth & environments).
	initSQLite()

	// Create a Gin router.
	router := gin.Default()

	// Auth endpoints.
	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

	// Environment endpoints:
	// - Only admin/superadmin can create new environments
	router.POST("/environments", AuthMiddleware(), AdminMiddleware(), createEnvironment)
	// - Any authenticated user can list them
	router.GET("/environments", AuthMiddleware(), listEnvironments)

	// Database/Collections routes by environment ID
	router.GET("/environments/:id/databases", AuthMiddleware(), getDatabasesForEnv)
	router.GET("/environments/:id/databases/:dbName/collections", AuthMiddleware(), getCollectionsForEnv)

	// Start the API server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// -------------------------------------------------------------------
// Auth-Related Handlers
// -------------------------------------------------------------------

// registerUser registers a new user in the local SQLite DB.
func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate required fields.
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}
	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Insert the new user into SQLite.
	stmt, err := db.Prepare(`INSERT INTO users(first_name, last_name, email, company, password, role)
	                         VALUES(?,?,?,?,?,?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	res, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Company, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last insert id"})
		return
	}
	user.ID = int(id)
	user.Password = "" // do not return the password
	c.JSON(http.StatusOK, gin.H{"user": user})
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

	// Retrieve the user from SQLite by email.
	var user User
	row := db.QueryRow(`SELECT id, first_name, last_name, email, company, password, role
	                    FROM users WHERE email = ?`, credentials.Email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Company, &user.Password, &user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create a JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(), // token expires in 72 hours
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

// AuthMiddleware extracts the JWT token from "Authorization: Bearer <token>",
// verifies it, and loads the user into the context.
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
			// Load user from DB
			var user User
			row := db.QueryRow(`SELECT id, first_name, last_name, email, company, password, role
			                    FROM users WHERE id = ?`, int(userID))
			err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
				&user.Company, &user.Password, &user.Role)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			// Clear password before storing in context
			user.Password = ""
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
	}
}

// AdminMiddleware ensures the user has role "admin" or "superadmin".
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		user, ok := userVal.(User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if user.Role != "admin" && user.Role != "superadmin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - Admins Only"})
			return
		}
		c.Next()
	}
}

// -------------------------------------------------------------------
// Environment CRUD
// -------------------------------------------------------------------

// createEnvironment allows admin/superadmin to register a new MongoDB environment.
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

	// Current user
	userVal, _ := c.Get("user")
	user := userVal.(User)

	// Insert into DB
	stmt, err := db.Prepare(`
		INSERT INTO environments (name, connection_string, created_by)
		VALUES (?,?,?)
	`)
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

// listEnvironments returns all registered Mongo environments.
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

// getDatabasesForEnv connects to the environment's MongoDB and lists all databases.
func getDatabasesForEnv(c *gin.Context) {
	envIDStr := c.Param("id")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Fetch environment from DB
	var e Environment
	row := db.QueryRow(`SELECT id, name, connection_string, created_by
	                    FROM environments WHERE id = ?`, envID)
	err = row.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Connect to the environment's MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(e.ConnectionString))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer mongoClient.Disconnect(ctx)

	result, err := mongoClient.ListDatabases(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// getCollectionsForEnv lists all collections in a specific DB of the given environment.
func getCollectionsForEnv(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")

	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Fetch environment
	var e Environment
	row := db.QueryRow(`SELECT id, name, connection_string, created_by
	                    FROM environments WHERE id = ?`, envID)
	err = row.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Connect to that Mongo instance
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(e.ConnectionString))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer mongoClient.Disconnect(ctx)

	mongoDB := mongoClient.Database(dbName)

	// Get collection names
	collections, err := mongoDB.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var output []gin.H

	// For each collection, gather stats
	for _, coll := range collections {
		stats := bson.M{}
		if err := mongoDB.RunCommand(ctx, bson.D{{Key: "collStats", Value: coll}}).Decode(&stats); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to get stats for %s: %v", coll, err),
			})
			return
		}

		indexes, err := mongoDB.Collection(coll).Indexes().List(ctx)
		var indexList []bson.M
		if err == nil {
			for indexes.Next(ctx) {
				var idx bson.M
				if decodeErr := indexes.Decode(&idx); decodeErr == nil {
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
