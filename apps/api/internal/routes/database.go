package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDatabaseRoutes(rg *gin.RouterGroup) {
	dbGroup := rg.Group("/environments/:id/databases")
	// Only authentication
	dbGroup.Use(middleware.AuthMiddleware())

	// Handlers themselves do the permission checks for read/write on the environment/database
	dbGroup.GET("", handlers.GetDatabases)
	dbGroup.POST("", handlers.CreateDatabase)
	dbGroup.GET("/:dbName", handlers.GetDatabaseDetails)
	dbGroup.PUT("/:dbName", handlers.EditDatabase)
	dbGroup.DELETE("/:dbName", handlers.DeleteDatabase)
}
