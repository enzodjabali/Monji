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
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:strongpassword@mongodb:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Routes
	router.GET("/databases", getDatabases)
	router.GET("/databases/:name/collections", getCollections)

	// Run API
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// GET /databases - List all databases
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

// GET /databases/:name/collections - List all collections and their stats
func getCollections(c *gin.Context) {
	dbName := c.Param("name")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := mongoClient.Database(dbName)

	// Get collection names
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var output []gin.H

	// Gather stats for each collection
	for _, coll := range collections {
		stats := bson.M{}
		err := db.RunCommand(ctx, bson.D{{Key: "collStats", Value: coll}}).Decode(&stats)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats for " + coll + ": " + err.Error()})
			return
		}

		indexes, err := db.Collection(coll).Indexes().List(ctx)
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
