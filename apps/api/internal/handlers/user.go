package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"monji/internal/database"
	"monji/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles creating a new user.
func CreateUser(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Company   string `json:"company"`
		Password  string `json:"password" binding:"required"`
		Role      string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the calling user is admin or superadmin
	callingUserRaw, _ := c.Get("user")
	callingUser := callingUserRaw.(models.User)

	// If caller is admin (NOT superadmin), then they cannot create a superadmin
	if callingUser.Role == "admin" && req.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin cannot create superadmin users"})
		return
	}

	// Hash the password before storing it.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	stmt, err := database.DB.Prepare(`
		INSERT INTO users(first_name, last_name, email, company, password, role)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	res, err := stmt.Exec(req.FirstName, req.LastName, req.Email, req.Company, string(hashedPassword), req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inserted ID"})
		return
	}

	// Return the newly created user (without the password).
	user := models.User{
		ID:        int(id),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Company:   req.Company,
		Role:      req.Role,
		Password:  "",
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser updates an existing user's details.
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Company   *string `json:"company"`
		Password  *string `json:"password"`
		Role      *string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check the current user’s role
	callingUserRaw, _ := c.Get("user")
	callingUser := callingUserRaw.(models.User)

	// We also need to check the role of the user we are updating:
	var existingRole string
	err = database.DB.QueryRow("SELECT role FROM users WHERE id = ?", id).Scan(&existingRole)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If calling user is admin and the target is superadmin => forbidden
	if callingUser.Role == "admin" && existingRole == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot edit superadmin user"})
		return
	}

	// If calling user is admin and they are trying to set the user’s role to superadmin => forbidden
	if callingUser.Role == "admin" && req.Role != nil && *req.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin cannot grant superadmin role"})
		return
	}

	// Build the update query dynamically.
	query := "UPDATE users SET "
	var params []interface{}
	var updates []string

	if req.FirstName != nil {
		updates = append(updates, "first_name = ?")
		params = append(params, *req.FirstName)
	}
	if req.LastName != nil {
		updates = append(updates, "last_name = ?")
		params = append(params, *req.LastName)
	}
	if req.Email != nil {
		updates = append(updates, "email = ?")
		params = append(params, *req.Email)
	}
	if req.Company != nil {
		updates = append(updates, "company = ?")
		params = append(params, *req.Company)
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updates = append(updates, "password = ?")
		params = append(params, string(hashedPassword))
	}
	if req.Role != nil {
		updates = append(updates, "role = ?")
		params = append(params, *req.Role)
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields provided to update"})
		return
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetch and return the updated user (omit the password).
	var user models.User
	row := database.DB.QueryRow("SELECT id, first_name, last_name, email, company, role FROM users WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// DeleteUser removes a user from the database.
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check the current user’s role
	callingUserRaw, _ := c.Get("user")
	callingUser := callingUserRaw.(models.User)

	// Check the role of the user we’re deleting
	var targetRole string
	err = database.DB.QueryRow("SELECT role FROM users WHERE id = ?", id).Scan(&targetRole)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If caller is admin and the target is superadmin => forbid
	if callingUser.Role == "admin" && targetRole == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete superadmin user"})
		return
	}

	res, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListUsers retrieves all users from the database.
func ListUsers(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, first_name, last_name, email, company, role FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUser retrieves a single user by ID.
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	row := database.DB.QueryRow("SELECT id, first_name, last_name, email, company, role FROM users WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
