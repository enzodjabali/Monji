package handlers

import (
	"net/http"

	"monji/internal/models"

	"github.com/gin-gonic/gin"
)

// WhoAmI returns all information about the current connected user.
// It expects the AuthMiddleware to have set the user in the context.
func WhoAmI(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	usr, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": usr})
}
