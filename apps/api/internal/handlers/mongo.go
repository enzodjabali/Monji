package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"monji/internal/database"
	"monji/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetDatabases lists databases for the specified MongoDB environment.
func GetDatabases(c *gin.Context) {
	envIDStr := c.Param("id")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
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

// GetCollections lists collections and basic stats for a given database.
func GetCollections(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
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
		output = append(output, gin.H{
			"name":           coll,
			"count":          stats["count"],
			"size":           stats["size"],
			"storageSize":    stats["storageSize"],
			"totalIndexSize": stats["totalIndexSize"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"database":    dbName,
		"collections": output,
	})
}

// GetDocuments retrieves all documents from a specific collection.
// (For production, consider adding pagination.)
func GetDocuments(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	collName := c.Param("collName")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection(collName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var documents []bson.M
	if err = cursor.All(ctx, &documents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode documents: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"database":   dbName,
		"collection": collName,
		"documents":  documents,
	})
}
