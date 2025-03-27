package handlers

import (
	"net/http"

	"monji/internal/models"

	"github.com/gin-gonic/gin"
)

// WhoAmI returns all information about the current authenticated user
// PLUS their environment & database permissions.
func WhoAmI(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	usr, ok := userRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}

	// Fetch all permissions for this user
	perms, err := fetchUserPermissions(usr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":        usr,
		"permissions": perms,
	})
}
