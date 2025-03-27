package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDocumentRoutes(rg *gin.RouterGroup) {
	docGroup := rg.Group("/environments/:id/databases/:dbName/collections/:collName/documents")
	docGroup.Use(middleware.AuthMiddleware())

	// The handler code checks read/write permission on the DB
	docGroup.GET("", handlers.GetDocuments)
	docGroup.POST("", handlers.CreateDocument)
	docGroup.PUT("/:docID", handlers.UpdateDocument)
	docGroup.DELETE("/:docID", handlers.DeleteDocument)
	docGroup.GET("/:docID", handlers.GetDocument)
}
