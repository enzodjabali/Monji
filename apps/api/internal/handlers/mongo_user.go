package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"monji/internal/database"
	"monji/internal/models"
	"monji/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateMongoUser creates a new MongoDB user on the target database.
func CreateMongoUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not in context"})
		return
	}
	currentUser, ok := currentUserRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}
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
	if isAtlas(env.ConnectionString) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Managing MongoDB users is not available on Atlas environments"})
		return
	}
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBWrite, err := utils.HasDBPermission(currentUser, envID, dbName, "write")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBWrite {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to create Mongo users on this DB"})
			return
		}
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Roles    []struct {
			Role string `json:"role"`
			Db   string `json:"db"`
		} `json:"roles"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Username == "" || req.Password == "" || len(req.Roles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, password and roles are required"})
		return
	}
	var rolesArr []bson.M
	for _, r := range req.Roles {
		if r.Role == "" || r.Db == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Each role must have both role and db defined"})
			return
		}
		rolesArr = append(rolesArr, bson.M{"role": r.Role, "db": r.Db})
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
	command := bson.D{
		{"createUser", req.Username},
		{"pwd", req.Password},
		{"roles", rolesArr},
	}
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, command).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "result": result})
}

// ListMongoUsers lists all MongoDB users on the given database.
func ListMongoUsers(c *gin.Context) {
	currentUserRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not in context"})
		return
	}
	currentUser, ok := currentUserRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}
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
	if isAtlas(env.ConnectionString) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Managing MongoDB users is not available on Atlas environments"})
		return
	}
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBRead, err := utils.HasDBPermission(currentUser, envID, dbName, "read")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBRead {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to list Mongo users on this DB"})
			return
		}
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
	command := bson.D{{"usersInfo", 1}}
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, command).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetMongoUser fetches details for a specific MongoDB user.
func GetMongoUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not in context"})
		return
	}
	currentUser, ok := currentUserRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	username := c.Param("username")
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
	if isAtlas(env.ConnectionString) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Managing MongoDB users is not available on Atlas environments"})
		return
	}
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBRead, err := utils.HasDBPermission(currentUser, envID, dbName, "read")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBRead {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to read Mongo user info on this DB"})
			return
		}
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
	command := bson.D{{"usersInfo", username}}
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, command).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// EditMongoUser updates an existing MongoDB user's password and/or roles.
func EditMongoUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not in context"})
		return
	}
	currentUser, ok := currentUserRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	username := c.Param("username")
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
	if isAtlas(env.ConnectionString) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Managing MongoDB users is not available on Atlas environments"})
		return
	}
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBWrite, err := utils.HasDBPermission(currentUser, envID, dbName, "write")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBWrite {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to edit Mongo users in this DB"})
			return
		}
	}
	var req struct {
		Password *string `json:"password,omitempty"`
		Roles    *[]struct {
			Role string `json:"role"`
			Db   string `json:"db"`
		} `json:"roles,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	decryptedConn, err := decrypt(env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}
	client, err := database.ConnectMongo(ctx, decryptedConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)
	updateDoc := bson.D{}
	if req.Password != nil {
		updateDoc = append(updateDoc, bson.E{"pwd", *req.Password})
	}
	if req.Roles != nil {
		var rolesArr []bson.M
		for _, r := range *req.Roles {
			if r.Role == "" || r.Db == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Each role must have both role and db defined"})
				return
			}
			rolesArr = append(rolesArr, bson.M{"role": r.Role, "db": r.Db})
		}
		updateDoc = append(updateDoc, bson.E{"roles", rolesArr})
	}
	if len(updateDoc) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update parameters provided"})
		return
	}
	command := bson.D{{"updateUser", username}}
	for _, e := range updateDoc {
		command = append(command, e)
	}
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, command).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "result": result})
}

// DeleteMongoUser removes a MongoDB user from the given database.
func DeleteMongoUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not in context"})
		return
	}
	currentUser, ok := currentUserRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	username := c.Param("username")
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
	if isAtlas(env.ConnectionString) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Managing MongoDB users is not available on Atlas environments"})
		return
	}
	isAdmin := utils.IsAdmin(currentUser)
	if !isAdmin {
		hasDBWrite, err := utils.HasDBPermission(currentUser, envID, dbName, "write")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasDBWrite {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permission to delete Mongo users in this DB"})
			return
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	decryptedConn, err := decrypt(env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt connection string: " + err.Error()})
		return
	}
	client, err := database.ConnectMongo(ctx, decryptedConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)
	command := bson.D{{"dropUser", username}}
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, command).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to drop user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "result": result})
}

// isAtlas is a helper that returns true if the connection string likely points to an Atlas environment.
func isAtlas(connString string) bool {
	return strings.Contains(connString, "mongodb.net")
}
