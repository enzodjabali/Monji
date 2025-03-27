package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"monji/internal/database"
)

// SetUserEnvironmentPermission sets or updates a user's permission on a given environment.
// Endpoint: POST /users/:userId/environments/:envId/permissions
// Body: { "permission": "readOnly" } or "readAndWrite" or "none"
func SetUserEnvironmentPermission(c *gin.Context) {
	// user must be admin or superadmin (enforced by AdminMiddleware)
	userIdStr := c.Param("userId")
	envIdStr := c.Param("envId")

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}
	envID, err := strconv.Atoi(envIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid envId"})
		return
	}

	var body struct {
		Permission string `json:"permission"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate permission
	if body.Permission != "none" && body.Permission != "readOnly" && body.Permission != "readAndWrite" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission (use 'none', 'readOnly' or 'readAndWrite')"})
		return
	}

	// If permission == "none", we can just remove the row from user_env_permissions
	if body.Permission == "none" {
		res, err := database.DB.Exec(
			`DELETE FROM user_env_permissions WHERE user_id = ? AND environment_id = ?`,
			userID, envID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rows, _ := res.RowsAffected()
		c.JSON(http.StatusOK, gin.H{
			"message":       "Environment permission removed",
			"rowsAffected":  rows,
			"user_id":       userID,
			"environmentId": envID,
		})
		return
	}

	// Otherwise, upsert the row
	stmt, err := database.DB.Prepare(`
		INSERT INTO user_env_permissions (user_id, environment_id, permission)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, environment_id)
		DO UPDATE SET permission=excluded.permission
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare upsert statement"})
		return
	}
	_, err = stmt.Exec(userID, envID, body.Permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert environment permission: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Environment permission set successfully",
		"user_id":       userID,
		"environmentId": envID,
		"permission":    body.Permission,
	})
}

// SetUserDBPermission sets or updates a user's permission on a specific database in a given environment.
// Endpoint: POST /users/:userId/environments/:envId/databases/:dbName/permissions
// Body: { "permission": "readOnly" } or "readAndWrite" or "none"
func SetUserDBPermission(c *gin.Context) {
	userIdStr := c.Param("userId")
	envIdStr := c.Param("envId")
	dbName := c.Param("dbName")

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}
	envID, err := strconv.Atoi(envIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid envId"})
		return
	}

	var body struct {
		Permission string `json:"permission"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate permission
	if body.Permission != "none" && body.Permission != "readOnly" && body.Permission != "readAndWrite" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission (use 'none', 'readOnly' or 'readAndWrite')"})
		return
	}

	// If permission == "none", remove row
	if body.Permission == "none" {
		res, err := database.DB.Exec(
			`DELETE FROM user_db_permissions WHERE user_id = ? AND environment_id = ? AND db_name = ?`,
			userID, envID, dbName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rows, _ := res.RowsAffected()
		c.JSON(http.StatusOK, gin.H{
			"message":       "DB permission removed",
			"rowsAffected":  rows,
			"user_id":       userID,
			"environmentId": envID,
			"dbName":        dbName,
		})
		return
	}

	// Otherwise, upsert the row
	stmt, err := database.DB.Prepare(`
		INSERT INTO user_db_permissions (user_id, environment_id, db_name, permission)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, environment_id, db_name)
		DO UPDATE SET permission=excluded.permission
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare upsert statement"})
		return
	}
	_, err = stmt.Exec(userID, envID, dbName, body.Permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert DB permission: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Database permission set successfully",
		"user_id":       userID,
		"environmentId": envID,
		"dbName":        dbName,
		"permission":    body.Permission,
	})
}
