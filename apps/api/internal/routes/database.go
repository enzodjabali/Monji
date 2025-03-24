package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDatabaseRoutes(rg *gin.RouterGroup) {
	dbGroup := rg.Group("/environments/:id/databases")
	{
		dbGroup.GET("", handlers.GetDatabases)
		dbGroup.POST("", middleware.AdminMiddleware(), handlers.CreateDatabase)
		dbGroup.PUT("/:dbName", middleware.AdminMiddleware(), handlers.EditDatabase)
		dbGroup.DELETE("/:dbName", middleware.AdminMiddleware(), handlers.DeleteDatabase)
	}
}
