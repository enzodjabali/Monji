package routes

import (
	"monji/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterDocumentRoutes(rg *gin.RouterGroup) {
	docGroup := rg.Group("/environments/:id/databases/:dbName/collections/:collName/documents")
	{
		docGroup.GET("", handlers.GetDocuments)
		docGroup.POST("", handlers.CreateDocument)
		docGroup.PUT("/:docID", handlers.UpdateDocument)
		docGroup.DELETE("/:docID", handlers.DeleteDocument)
	}
}
