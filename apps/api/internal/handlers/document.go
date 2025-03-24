package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"monji/internal/database"
	"monji/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetDocuments retrieves all documents from a collection.
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
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
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

	cursor, err := client.Database(dbName).Collection(collName).Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var documents []bson.M
	if err := cursor.All(ctx, &documents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode documents: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"database":   dbName,
		"collection": collName,
		"documents":  documents,
	})
}

// CreateDocument inserts a new document into a collection.
func CreateDocument(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	collName := c.Param("collName")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	var env models.Environment
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	var doc bson.M
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	res, err := client.Database(dbName).Collection(collName).InsertOne(ctx, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert document: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "Document created successfully",
		"insertedId": res.InsertedID,
	})
}

// UpdateDocument updates a document by its ObjectID.
func UpdateDocument(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	collName := c.Param("collName")
	docIDStr := c.Param("docID")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	var env models.Environment
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	docID, err := primitive.ObjectIDFromHex(docIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	res, err := client.Database(dbName).Collection(collName).UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{"$set": updateData},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "Document updated successfully",
		"matchedCount":  res.MatchedCount,
		"modifiedCount": res.ModifiedCount,
	})
}

// DeleteDocument deletes a document by its ObjectID.
func DeleteDocument(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	collName := c.Param("collName")
	docIDStr := c.Param("docID")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}
	var env models.Environment
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	docID, err := primitive.ObjectIDFromHex(docIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
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

	res, err := client.Database(dbName).Collection(collName).DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "Document deleted successfully",
		"deletedCount": res.DeletedCount,
	})
}
