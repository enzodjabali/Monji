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

		dbGroup.GET("", handlers.GetDatabases)
		dbGroup.POST("", handlers.CreateDatabase)
		dbGroup.PUT("/:dbName", handlers.EditDatabase)
		dbGroup.DELETE("/:dbName", handlers.DeleteDatabase)
	}
}
