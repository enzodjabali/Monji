package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"monji/internal/database"
	"monji/internal/models"
	"monji/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetCollections lists collections (with basic stats) in a database.
func GetCollections(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")

	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Check read permission on this DB
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
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to read this database"})
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

	collNames, err := client.Database(dbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list collections: " + err.Error()})
		return
	}

	var collections []gin.H
	for _, coll := range collNames {
		var stats bson.M
		if err := client.Database(dbName).RunCommand(ctx, bson.D{{Key: "collStats", Value: coll}}).Decode(&stats); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get stats for %s: %v", coll, err)})
			return
		}
		collections = append(collections, gin.H{
			"name":           coll,
			"count":          stats["count"],
			"size":           stats["size"],
			"storageSize":    stats["storageSize"],
			"totalIndexSize": stats["totalIndexSize"],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"database":    dbName,
		"collections": collections,
	})
}

// CreateCollection creates a new collection in a database.
func CreateCollection(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")

	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Check write permission on this DB
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
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to write in this database"})
			return
		}
	}

	var req struct {
		CollectionName string `json:"collectionName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.CollectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectionName is required"})
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

	collNames, err := client.Database(dbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list collections: " + err.Error()})
		return
	}
	for _, name := range collNames {
		if name == req.CollectionName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Collection already exists"})
			return
		}
	}

	if err := client.Database(dbName).CreateCollection(ctx, req.CollectionName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Collection created successfully",
		"database":   dbName,
		"collection": req.CollectionName,
	})
}

// EditCollection renames an existing collection.
func EditCollection(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	oldCollName := c.Param("collName")

	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Load environment
	var env models.Environment
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Check write permission
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
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to write in this database"})
			return
		}
	}

	var req struct {
		NewCollectionName string `json:"newCollectionName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NewCollectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newCollectionName is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	oldNamespace := fmt.Sprintf("%s.%s", dbName, oldCollName)
	newNamespace := fmt.Sprintf("%s.%s", dbName, req.NewCollectionName)
	cmd := bson.D{
		{"renameCollection", oldNamespace},
		{"to", newNamespace},
		{"dropTarget", false},
	}
	if err := client.Database("admin").RunCommand(ctx, cmd).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to rename collection: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Collection renamed successfully",
		"oldCollection": oldCollName,
		"newCollection": req.NewCollectionName,
	})
}

// DeleteCollection drops a collection from a database.
func DeleteCollection(c *gin.Context) {
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
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Check write permission
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
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to write in this database"})
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

	if err := client.Database(dbName).Collection(collName).Drop(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to drop collection: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Collection deleted successfully",
		"database":   dbName,
		"collection": collName,
	})
}

// GetCollectionDetails retrieves detailed info about a collection.
func GetCollectionDetails(c *gin.Context) {
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
	row := database.DB.QueryRow(`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Check read permission
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
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to read this collection"})
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

	var stats bson.M
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{Key: "collStats", Value: collName}}).Decode(&stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get collection stats: " + err.Error()})
		return
	}

	cursor, err := client.Database(dbName).Collection(collName).Indexes().List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list indexes: " + err.Error()})
		return
	}
	var indexes []bson.M
	for cursor.Next(ctx) {
		var idx bson.M
		if err := cursor.Decode(&idx); err == nil {
			indexes = append(indexes, idx)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"database":   dbName,
		"collection": collName,
		"stats":      stats,
		"indexes":    indexes,
	})
}
