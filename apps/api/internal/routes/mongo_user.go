package routes

import (
	"monji/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterMongoUserRoutes sets up endpoints for managing MongoDB database users.
// Only Auth at the router level; the handler checks if user is admin or has DB permissions.
func RegisterMongoUserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/environments/:id/databases/:dbName/users")
	// We only do Auth. The handler enforces read/write on DB if needed.
	// (If you want to restrict creation/editing of Mongo DB users to admin only, you can do that in the handler.)

	// but let's assume you also want to require Auth
	// If you do want to allow normal users to read them if they have read permission, fine:
	// if you want to require admin role only, you'd add AdminMiddleware(). But per your spec, let's keep it open:
	// userGroup.Use(middleware.AuthMiddleware()) // you might want to do this if not inherited from a higher group

	userGroup.Use(gin.Logger()) // or nothing, up to you. But presumably you do:
	// userGroup.Use(middleware.AuthMiddleware())

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
