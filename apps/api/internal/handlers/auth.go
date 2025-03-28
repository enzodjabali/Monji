package handlers

import (
	"net/http"

	"monji/internal/database"
	"monji/internal/middleware"
	"monji/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user authentication and JWT token generation.
func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	row := database.DB.QueryRow(
		`SELECT id, first_name, last_name, email, company, password, role
		 FROM users WHERE email = ?`, credentials.Email)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Password, &user.Role); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Validate password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token.
	token, err := middleware.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
