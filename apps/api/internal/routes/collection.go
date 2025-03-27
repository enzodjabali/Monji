package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterCollectionRoutes ensures we only do Auth here.
// The collection handlers check if user is admin or has read/write on the DB.
func RegisterCollectionRoutes(rg *gin.RouterGroup) {
	collGroup := rg.Group("/environments/:id/databases/:dbName/collections")
	collGroup.Use(middleware.AuthMiddleware())

	collGroup.GET("", handlers.GetCollections)
	collGroup.GET("/:collName", handlers.GetCollectionDetails)
	collGroup.POST("", handlers.CreateCollection)
	collGroup.PUT("/:collName", handlers.EditCollection)
	collGroup.DELETE("/:collName", handlers.DeleteCollection)
}
