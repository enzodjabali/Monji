package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterDocumentRoutes ensures both Auth and Admin for all document endpoints.
func RegisterDocumentRoutes(rg *gin.RouterGroup) {
	docGroup := rg.Group("/environments/:id/databases/:dbName/collections/:collName/documents")
	{
		docGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

		docGroup.GET("", handlers.GetDocuments)
		docGroup.POST("", handlers.CreateDocument)
		docGroup.PUT("/:docID", handlers.UpdateDocument)
		docGroup.DELETE("/:docID", handlers.DeleteDocument)
	}
}
