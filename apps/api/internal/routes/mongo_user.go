package routes

import (
	"monji/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterMongoUserRoutes sets up endpoints for managing MongoDB database users.
func RegisterMongoUserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/environments/:id/databases/:dbName/users")
	{
		// Create a new MongoDB user.
		userGroup.POST("", handlers.CreateMongoUser)
		// List all users for the database.
		userGroup.GET("", handlers.ListMongoUsers)
		// Get details of a specific user.
		userGroup.GET("/:username", handlers.GetMongoUser)
		// Update (edit) an existing user.
		userGroup.PUT("/:username", handlers.EditMongoUser)
		// Delete a user.
		userGroup.DELETE("/:username", handlers.DeleteMongoUser)
	}
}
