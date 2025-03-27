package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"monji/internal/database"
	"monji/internal/models"
	"monji/internal/utils"

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

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by
		   FROM environments
		  WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Check read permission on DB
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBRead, err := utils.HasDBPermission(currentUser, envID, dbName, "read")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBRead {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to read documents in this DB"})
			return
		}
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

// GetDocument fetches a single document from a collection by its _id (hex or string).
func GetDocument(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	collName := c.Param("collName")
	docIDStr := c.Param("docID")

	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by
		   FROM environments
		  WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		// Check read
		hasDBRead, err := utils.HasDBPermission(currentUser, envID, dbName, "read")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBRead {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to read documents in this DB"})
			return
		}
	}

	// Try converting docIDStr to an ObjectID. If it fails, treat it as a string _id.
	var filter bson.M
	objID, err := primitive.ObjectIDFromHex(docIDStr)
	if err != nil {
		// Not a valid hex => treat as string
		filter = bson.M{"_id": docIDStr}
	} else {
		// Valid hex => treat as ObjectID
		filter = bson.M{"_id": objID}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	var result bson.M
	if err := client.Database(dbName).Collection(collName).FindOne(ctx, filter).Decode(&result); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"database":   dbName,
		"collection": collName,
		"document":   result,
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

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by
		   FROM environments
		  WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Check write
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBWrite, err := utils.HasDBPermission(currentUser, envID, dbName, "write")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBWrite {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to write documents in this DB"})
			return
		}
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

// UpdateDocument updates a document by its _id (hex or string).
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

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by
		   FROM environments
		  WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Attempt to parse docIDStr as ObjectID; fallback to string
	var filter bson.M
	objID, err := primitive.ObjectIDFromHex(docIDStr)
	if err != nil {
		filter = bson.M{"_id": docIDStr}
	} else {
		filter = bson.M{"_id": objID}
	}

	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove _id field if present to avoid immutable field error
	delete(updateData, "_id")

	// Check write
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBWrite, err := utils.HasDBPermission(currentUser, envID, dbName, "write")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBWrite {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to write documents in this DB"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	res, err := client.Database(dbName).
		Collection(collName).
		UpdateOne(ctx, filter, bson.M{"$set": updateData})
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

// DeleteDocument deletes a document by its _id (hex or string).
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

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by
		   FROM environments
		  WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Attempt to parse docIDStr as ObjectID; fallback to string
	var filter bson.M
	objID, err := primitive.ObjectIDFromHex(docIDStr)
	if err != nil {
		filter = bson.M{"_id": docIDStr}
	} else {
		filter = bson.M{"_id": objID}
	}

	// Check write
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBWrite, err := utils.HasDBPermission(currentUser, envID, dbName, "write")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBWrite {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to delete documents in this DB"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	res, err := client.Database(dbName).Collection(collName).DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Document deleted successfully",
		"deletedCount": res.DeletedCount,
	})
}
