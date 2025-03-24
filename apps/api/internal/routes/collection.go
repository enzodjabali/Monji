package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCollectionRoutes(rg *gin.RouterGroup) {
	collGroup := rg.Group("/environments/:id/databases/:dbName/collections")
	{
		collGroup.GET("", handlers.GetCollections)
		// Getting detailed info requires admin privileges.
		collGroup.GET("/:collName", middleware.AdminMiddleware(), handlers.GetCollectionDetails)
		collGroup.POST("", middleware.AdminMiddleware(), handlers.CreateCollection)
		collGroup.PUT("/:collName", middleware.AdminMiddleware(), handlers.EditCollection)
		collGroup.DELETE("/:collName", middleware.AdminMiddleware(), handlers.DeleteCollection)
	}
}
