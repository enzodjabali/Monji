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

// GetDatabases lists Mongo databases in the specified environment.
// It decrypts the stored connection string before connecting.
func GetDatabases(c *gin.Context) {
	envIDStr := c.Param("id")
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

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	isAdmin := utils.IsAdmin(currentUser)

	if !isAdmin {
		hasEnvRead, err := utils.HasEnvPermission(currentUser, envID, "read")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasEnvRead {
			c.JSON(http.StatusOK, gin.H{
				"Databases": []interface{}{},
				"TotalSize": 0,
			})
			return
		}
	}

	decryptedConn, err := decrypt(env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, decryptedConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	dbs, err := client.ListDatabases(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resultList []map[string]interface{}
	var totalSize float64

	for _, dbInfo := range dbs.Databases {
		if isAdmin {
			resultList = append(resultList, map[string]interface{}{
				"Name":         dbInfo.Name,
				"SizeOnDisk":   dbInfo.SizeOnDisk,
				"Empty":        dbInfo.Empty,
				"myPermission": "readAndWrite",
			})
			totalSize += float64(dbInfo.SizeOnDisk)
		} else {
			hasDbRead, err := utils.HasDBPermission(currentUser, envID, dbInfo.Name, "read")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if hasDbRead {
				perm := getDbPermissionString(currentUser, envID, dbInfo.Name)
				resultList = append(resultList, map[string]interface{}{
					"Name":         dbInfo.Name,
					"SizeOnDisk":   dbInfo.SizeOnDisk,
					"Empty":        dbInfo.Empty,
					"myPermission": perm,
				})
				totalSize += float64(dbInfo.SizeOnDisk)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Databases": resultList,
		"TotalSize": totalSize,
	})
}

// CreateDatabase creates a new Mongo database by creating an initial collection.
func CreateDatabase(c *gin.Context) {
	envIDStr := c.Param("id")
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

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	hasEnvWrite, err := utils.HasEnvPermission(currentUser, envID, "write")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasEnvWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission on this environment"})
		return
	}

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

	decryptedConn, err := decrypt(env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, decryptedConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	dbList, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list databases: " + err.Error()})
		return
	}
	for _, existingName := range dbList {
		if existingName == req.DbName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Database already exists"})
			return
		}
	}

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

// EditDatabase renames a database by moving all collections to a new database.
func EditDatabase(c *gin.Context) {
	envIDStr := c.Param("id")
	oldDbName := c.Param("dbName")

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

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	hasEnvWrite, err := utils.HasEnvPermission(currentUser, envID, "write")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasEnvWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission on environment"})
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

	decryptedConn, err := decrypt(env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := database.ConnectMongo(ctx, decryptedConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

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

	collNames, err := client.Database(oldDbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list collections: " + err.Error()})
		return
	}

	for _, coll := range collNames {
		oldNamespace := fmt.Sprintf("%s.%s", oldDbName, coll)
		newNamespace := fmt.Sprintf("%s.%s", req.NewDbName, coll)
		cmd := bson.D{
			{Key: "renameCollection", Value: oldNamespace},
			{Key: "to", Value: newNamespace},
			{Key: "dropTarget", Value: false},
		}
		if err := client.Database("admin").RunCommand(ctx, cmd).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to rename collection %s: %v", coll, err)})
			return
		}
	}

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

// DeleteDatabase drops a Mongo database.
func DeleteDatabase(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")

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

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	hasEnvWrite, err := utils.HasEnvPermission(currentUser, envID, "write")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasEnvWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission on environment"})
		return
	}

	decryptedConn, err := decrypt(env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, decryptedConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}
	defer client.Disconnect(ctx)

	if err := client.Database(dbName).Drop(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to drop database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Database deleted successfully",
		"database": dbName,
	})
}

// GetDatabaseDetails returns detailed info about a specific MongoDB database,
// including statistics, collections, and the current user's DB-level permission.
func GetDatabaseDetails(c *gin.Context) {
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")

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

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)

	hasEnvRead, err := utils.HasEnvPermission(currentUser, envID, "read")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasEnvRead {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission on environment"})
		return
	}

	hasDbRead, err := utils.HasDBPermission(currentUser, envID, dbName, "read")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasDbRead {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission on database"})
		return
	}

	decryptedConn, err := decrypt(env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, decryptedConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	var stats bson.M
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database stats: " + err.Error()})
		return
	}

	collNames, err := client.Database(dbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list collections: " + err.Error()})
		return
	}

	myPerm := "readAndWrite"
	if !utils.IsAdmin(currentUser) {
		myPerm = getDbPermissionString(currentUser, envID, dbName)
	}

	c.JSON(http.StatusOK, gin.H{
		"database":     dbName,
		"stats":        stats,
		"collections":  collNames,
		"myPermission": myPerm,
	})
}
