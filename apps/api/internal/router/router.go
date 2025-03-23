package router

import (
	"monji/internal/config"
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter registers all routes and returns a Gin engine.
func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Public route.
	router.POST("/login", handlers.Login)

	// Protected routes.
	api := router.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		// Environment endpoints.
		api.GET("/environments", handlers.ListEnvironments)
		api.GET("/environments/:id", handlers.GetEnvironment)
		api.POST("/environments", middleware.AdminMiddleware(), handlers.CreateEnvironment)
		api.PUT("/environments/:id", middleware.AdminMiddleware(), handlers.UpdateEnvironment)
		api.DELETE("/environments/:id", middleware.AdminMiddleware(), handlers.DeleteEnvironment)

		// Mongo queries.
		api.GET("/environments/:id/databases", handlers.GetDatabases)
		api.GET("/environments/:id/databases/:dbName/collections", handlers.GetCollections)
		api.GET("/environments/:id/databases/:dbName/collections/:collName/documents", handlers.GetDocuments)
	}

	return router
}
