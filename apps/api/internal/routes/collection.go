package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterCollectionRoutes ensures both Auth and Admin for all collection endpoints.
func RegisterCollectionRoutes(rg *gin.RouterGroup) {
	collGroup := rg.Group("/environments/:id/databases/:dbName/collections")
	{
		collGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

		collGroup.GET("", handlers.GetCollections)
		collGroup.GET("/:collName", handlers.GetCollectionDetails)
		collGroup.POST("", handlers.CreateCollection)
		collGroup.PUT("/:collName", handlers.EditCollection)
		collGroup.DELETE("/:collName", handlers.DeleteCollection)
	}
}
