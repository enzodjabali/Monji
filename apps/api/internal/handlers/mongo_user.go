package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"monji/internal/database"
	"monji/internal/models"
)

// CreateMongoUser creates a new MongoDB user on the target database.
func CreateMongoUser(c *gin.Context) {
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

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		// Roles is an array of roles with the role name and target database.
		Roles []struct {
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

	// Prepare roles array for the createUser command.
	var rolesArr []bson.M
	for _, r := range req.Roles {
		if r.Role == "" || r.Db == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Each role must have both role and db defined"})
			return
		}
		rolesArr = append(rolesArr, bson.M{"role": r.Role, "db": r.Db})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	// Run the createUser command on the target database.
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	// Use the usersInfo command to list users.
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
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	username := c.Param("username")

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	// Use the usersInfo command to get details for the specific user.
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
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	username := c.Param("username")

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
		Password *string `json:"password,omitempty"`
		// Roles is optional; if provided, it will replace the existing roles.
		Roles *[]struct {
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
	client, err := database.ConnectMongo(ctx, env.ConnectionString)
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

	// Use the updateUser command to change the user's password and/or roles.
	command := bson.D{{"updateUser", username}}
	for _, e := range updateDoc {
		command = append(command, bson.E{Key: e.Key, Value: e.Value})
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
	envIDStr := c.Param("id")
	dbName := c.Param("dbName")
	username := c.Param("username")

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB: " + err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	// Run the dropUser command to remove the user.
	command := bson.D{{"dropUser", username}}
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, command).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to drop user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "result": result})
}
