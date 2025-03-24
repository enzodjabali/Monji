package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterDatabaseRoutes ensures both Auth and Admin are required for all DB endpoints.
func RegisterDatabaseRoutes(rg *gin.RouterGroup) {
	dbGroup := rg.Group("/environments/:id/databases")
	{
		dbGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

		// List all databases.
		dbGroup.GET("", handlers.GetDatabases)
		// Create a new database.
		dbGroup.POST("", handlers.CreateDatabase)
		// Get details for a specific database.
		dbGroup.GET("/:dbName", handlers.GetDatabaseDetails)
		// Rename a database.
		dbGroup.PUT("/:dbName", handlers.EditDatabase)
		// Delete a database.
		dbGroup.DELETE("/:dbName", handlers.DeleteDatabase)
	}
}
