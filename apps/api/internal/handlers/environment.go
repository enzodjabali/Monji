package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"monji/internal/database"
	"monji/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateEnvironment creates a new MongoDB environment configuration.
func CreateEnvironment(c *gin.Context) {
	var req struct {
		Name             string `json:"name"`
		ConnectionString string `json:"connection_string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" || req.ConnectionString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing environment name or connection string"})
		return
	}

	// Retrieve current user from context.
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	usr := user.(models.User)

	stmt, err := database.DB.Prepare(`INSERT INTO environments (name, connection_string, created_by) VALUES (?,?,?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	res, err := stmt.Exec(req.Name, req.ConnectionString, usr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last insert id"})
		return
	}

	env := models.Environment{
		ID:               int(id),
		Name:             req.Name,
		ConnectionString: req.ConnectionString,
		CreatedBy:        usr.ID,
	}
	c.JSON(http.StatusOK, gin.H{"environment": env})
}

// ListEnvironments returns all stored MongoDB environment configurations.
func ListEnvironments(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, name, connection_string, created_by FROM environments`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var envs []models.Environment
	for rows.Next() {
		var e models.Environment
		if err := rows.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		envs = append(envs, e)
	}
	c.JSON(http.StatusOK, gin.H{"environments": envs})
}

// UpdateEnvironment edits an existing environment configuration.
// It accepts a JSON payload with the fields to update.
func UpdateEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	var req struct {
		Name             string `json:"name"`
		ConnectionString string `json:"connection_string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" && req.ConnectionString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update parameters provided"})
		return
	}

	// Build dynamic update query.
	query := "UPDATE environments SET "
	var params []interface{}
	if req.Name != "" {
		query += "name = ?"
		params = append(params, req.Name)
	}
	if req.ConnectionString != "" {
		if len(params) > 0 {
			query += ", "
		}
		query += "connection_string = ?"
		params = append(params, req.ConnectionString)
	}
	query += " WHERE id = ?"
	params = append(params, id)

	res, err := database.DB.Exec(query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// Query the updated environment.
	var env models.Environment
	row := database.DB.QueryRow("SELECT id, name, connection_string, created_by FROM environments WHERE id = ?", id)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"environment": env})
}

// DeleteEnvironment removes an environment configuration.
func DeleteEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	res, err := database.DB.Exec("DELETE FROM environments WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Environment deleted successfully"})
}
