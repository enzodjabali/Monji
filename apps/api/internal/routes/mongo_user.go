package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterMongoUserRoutes sets up endpoints for managing MongoDB database users.
func RegisterMongoUserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/environments/:id/databases/:dbName/users")

	// Make sure to use AuthMiddleware() so "user" is set in the context.
	userGroup.Use(middleware.AuthMiddleware())

	userGroup.POST("", handlers.CreateMongoUser)
	userGroup.GET("", handlers.ListMongoUsers)
	userGroup.GET("/:username", handlers.GetMongoUser)
	userGroup.PUT("/:username", handlers.EditMongoUser)
	userGroup.DELETE("/:username", handlers.DeleteMongoUser)
}
