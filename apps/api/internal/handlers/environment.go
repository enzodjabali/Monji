package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"monji/internal/database"
	"monji/internal/models"
	"monji/internal/utils"

	"github.com/gin-gonic/gin"
)

// getEnvPermissionString returns the environment-level permission for the given user.
// If user is admin/superadmin => "readAndWrite".
// Otherwise, we look in user_env_permissions for "readOnly" / "readAndWrite" / no row => "none".
func getEnvPermissionString(user models.User, envID int) string {
	if utils.IsAdmin(user) {
		return "readAndWrite"
	}

	row := database.DB.QueryRow(
		`SELECT permission FROM user_env_permissions WHERE user_id = ? AND environment_id = ?`,
		user.ID, envID,
	)
	var perm string
	err := row.Scan(&perm)
	if err != nil {
		// if no row or error => treat as "none"
		return "none"
	}
	return perm
}

// CreateEnvironment creates a new MongoDB environment configuration.
func CreateEnvironment(c *gin.Context) {
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)
	if !utils.IsAdmin(currentUser) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin or superadmin can create environments"})
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
	if req.Name == "" || req.ConnectionString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing environment name or connection string"})
		return
	}

	stmt, err := database.DB.Prepare(`INSERT INTO environments (name, connection_string, created_by) VALUES (?,?,?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	res, err := stmt.Exec(req.Name, req.ConnectionString, currentUser.ID)
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
		CreatedBy:        currentUser.ID,
	}
	c.JSON(http.StatusOK, gin.H{"environment": env})
}

// ListEnvironments returns all environment configurations for admin/superadmin,
// or only the permitted environments for a normal user.
//
// Now also returns "myPermission" for each environment object.
func ListEnvironments(c *gin.Context) {
	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)

	if utils.IsAdmin(currentUser) {
		// admin => list all environments
		rows, err := database.DB.Query(`SELECT id, name, connection_string, created_by FROM environments`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var envs []gin.H
		for rows.Next() {
			var e models.Environment
			if err := rows.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy); err != nil {
				if err == sql.ErrNoRows {
					break
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// admin => "readAndWrite"
			res := gin.H{
				"id":                e.ID,
				"name":              e.Name,
				"connection_string": e.ConnectionString,
				"created_by":        e.CreatedBy,
				"myPermission":      "readAndWrite",
			}
			envs = append(envs, res)
		}
		c.JSON(http.StatusOK, gin.H{"environments": envs})
		return
	}

	// normal user => show only envs they have read permission for
	query := `
	SELECT e.id, e.name, e.connection_string, e.created_by,
	       p.permission
	  FROM environments e
	  JOIN user_env_permissions p ON e.id = p.environment_id
	 WHERE p.user_id = ?
	   AND (p.permission = 'readOnly' OR p.permission = 'readAndWrite');
	`
	rows, err := database.DB.Query(query, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var envs []gin.H
	for rows.Next() {
		var (
			e    models.Environment
			perm string
		)
		if err := rows.Scan(&e.ID, &e.Name, &e.ConnectionString, &e.CreatedBy, &perm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		res := gin.H{
			"id":                e.ID,
			"name":              e.Name,
			"connection_string": e.ConnectionString,
			"created_by":        e.CreatedBy,
			"myPermission":      perm,
		}
		envs = append(envs, res)
	}
	c.JSON(http.StatusOK, gin.H{"environments": envs})
}

// GetEnvironment fetches details for a single environment, respecting permissions,
// and adds "myPermission" at the top level.
func GetEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)

	var env models.Environment
	row := database.DB.QueryRow("SELECT id, name, connection_string, created_by FROM environments WHERE id = ?", id)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment not found"})
		return
	}

	// compute myPermission
	myPerm := getEnvPermissionString(currentUser, env.ID)
	if myPerm == "none" {
		// user has no read => either 403 or return a restricted response.
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to read this environment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"environment": gin.H{
			"id":                env.ID,
			"name":              env.Name,
			"connection_string": env.ConnectionString,
			"created_by":        env.CreatedBy,
		},
		"myPermission": myPerm,
	})
}

// UpdateEnvironment updates an environment configuration. (Unchanged from your original, no "myPermission" needed.)
func UpdateEnvironment(c *gin.Context) {
	// ... unchanged, no need to show "myPermission" for update
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)

	hasWrite, err := utils.HasEnvPermission(currentUser, id, "write")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to write this environment"})
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

	query := "UPDATE environments SET "
	var params []interface{}
	var updates []string

	if req.Name != "" {
		updates = append(updates, "name = ?")
		params = append(params, req.Name)
	}
	if req.ConnectionString != "" {
		updates = append(updates, "connection_string = ?")
		params = append(params, req.ConnectionString)
	}
	query += joinUpdates(updates, ", ") + " WHERE id = ?"
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

	var env models.Environment
	row := database.DB.QueryRow("SELECT id, name, connection_string, created_by FROM environments WHERE id = ?", id)
	if err := row.Scan(&env.ID, &env.Name, &env.ConnectionString, &env.CreatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"environment": env})
}

// DeleteEnvironment removes an environment configuration. (Unchanged, no "myPermission" needed.)
func DeleteEnvironment(c *gin.Context) {
	// ... unchanged from your original
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
		return
	}

	currentUserRaw, _ := c.Get("user")
	currentUser := currentUserRaw.(models.User)

	hasWrite, err := utils.HasEnvPermission(currentUser, id, "write")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to delete this environment"})
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
