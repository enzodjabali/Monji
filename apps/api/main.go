package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoClient is a global variable to hold the MongoDB client instance.
var mongoClient *mongo.Client

func main() {
	// Set up a connection to MongoDB using the connection string "mongodb://mongo:27017"
	// "mongo" is the hostname of the MongoDB container defined in docker-compose.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:strongpassword@mongodb:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create a Gin router instance.
	router := gin.Default()

	// Define an endpoint that returns the list of databases.
	router.GET("/databases", getDatabases)

	// Run the API server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// getDatabases handles GET /databases and returns the list of databases along with stats.
func getDatabases(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ListDatabases returns a structure with the Databases field containing details such as name, sizeOnDisk, and empty flag.
	result, err := mongoClient.ListDatabases(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
