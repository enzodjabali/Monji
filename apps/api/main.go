package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

// Global variables
var (
	mongoClient *mongo.Client
	db          *sql.DB                    // SQLite database connection
	jwtSecret   = []byte("supersecretkey") // Ideally, set this via an environment variable.
)

// User represents our user model.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Company   string `json:"company,omitempty"`
	Password  string `json:"password"` // <-- ALLOW JSON input
	Role      string `json:"role"`
}

func initSQLite() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	// Create the users table if it doesn't exist.
	createTableSQL := `
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
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
}

func main() {
	// Connect to MongoDB.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:strongpassword@mongodb:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize SQLite (for user auth).
	initSQLite()

	// Create a Gin router.
	router := gin.Default()

	// Auth endpoints.
	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

	// Existing MongoDB endpoints.
	router.GET("/databases", getDatabases)
	router.GET("/databases/:name/collections", getCollections)

	// Start the API server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// registerUser registers a new user.
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
	stmt, err := db.Prepare("INSERT INTO users(first_name, last_name, email, company, password, role) VALUES(?,?,?,?,?,?)")
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
	// Retrieve the user from SQLite.
	var user User
	row := db.QueryRow("SELECT id, first_name, last_name, email, company, password, role FROM users WHERE email = ?", credentials.Email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Password, &user.Role)
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

// getDatabases lists all databases in MongoDB.
func getDatabases(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := mongoClient.ListDatabases(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// getCollections lists all collections in a given MongoDB database along with their stats.
func getCollections(c *gin.Context) {
	dbName := c.Param("name")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoDB := mongoClient.Database(dbName)

	// Get collection names.
	collections, err := mongoDB.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var output []gin.H

	// For each collection, get statistics and indexes.
	for _, coll := range collections {
		stats := bson.M{}
		err := mongoDB.RunCommand(ctx, bson.D{{Key: "collStats", Value: coll}}).Decode(&stats)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats for " + coll + ": " + err.Error()})
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
