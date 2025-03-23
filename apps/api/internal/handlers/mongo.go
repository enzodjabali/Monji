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

// CreateDatabase creates a new database in the specified MongoDB environment
// by creating an initial collection. In MongoDB a database is created when
// its first collection is created.
func CreateDatabase(c *gin.Context) {
	envIDStr := c.Param("id")
	envID, err := strconv.Atoi(envIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	// Retrieve the environment from SQLite.
	var env models.Environment
	row := database.DB.QueryRow(
		`SELECT id, name, connection_string, created_by FROM environments WHERE id = ?`, envID)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Parse JSON payload.
	var req struct {
		DbName            string `json:"dbName"`
		InitialCollection string `json:"initialCollection"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.DbName == "" || req.InitialCollection == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both dbName and initialCollection are required"})
		return
	}

	// Connect to MongoDB using the environment's connection string.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	// Check if the database already exists.
	dbList, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list databases: " + err.Error()})
		return
	}
	for _, dbName := range dbList {
		if dbName == req.DbName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Database already exists"})
			return
		}
	}

	// Create the new database by creating the initial collection.
	if err := client.Database(req.DbName).CreateCollection(ctx, req.InitialCollection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Database created successfully",
		"database":          req.DbName,
		"initialCollection": req.InitialCollection,
	})
}

// EditDatabase renames a database by moving all collections from the old database
// to a new database. It expects a JSON payload with the field "newDbName".
// Note: MongoDB does not support renaming databases directly.
func EditDatabase(c *gin.Context) {
	envIDStr := c.Param("id")
	oldDbName := c.Param("dbName")
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

	var req struct {
		NewDbName string `json:"newDbName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NewDbName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newDbName is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	// Check if the target database already exists.
	dbList, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list databases: " + err.Error()})
		return
	}
	for _, name := range dbList {
		if name == req.NewDbName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Target database already exists"})
			return
		}
	}

	// List collections in the old database.
	collNames, err := client.Database(oldDbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list collections: " + err.Error()})
		return
	}

	// For each collection, move it to the new database using renameCollection.
	for _, coll := range collNames {
		oldNamespace := fmt.Sprintf("%s.%s", oldDbName, coll)
		newNamespace := fmt.Sprintf("%s.%s", req.NewDbName, coll)
		cmd := bson.D{
			{"renameCollection", oldNamespace},
			{"to", newNamespace},
			{"dropTarget", false},
		}
		// The renameCollection command must be run against the "admin" database.
		if err := client.Database("admin").RunCommand(ctx, cmd).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to rename collection %s: %v", coll, err)})
			return
		}
	}

	// Drop the old database (it should now be empty).
	if err := client.Database(oldDbName).Drop(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to drop old database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database renamed successfully",
		"oldName": oldDbName,
		"newName": req.NewDbName,
	})
}

// DeleteDatabase drops a database in the specified MongoDB environment.
func DeleteDatabase(c *gin.Context) {
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

	// Drop the database.
	if err := client.Database(dbName).Drop(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to drop database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Database deleted successfully",
		"database": dbName,
	})
}
