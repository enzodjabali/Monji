package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up the endpoints for user CRUD operations.
// These endpoints are protected by Auth and Admin middleware.
func RegisterUserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

	userGroup.POST("", handlers.CreateUser)
	userGroup.PUT("/:id", handlers.UpdateUser)
	userGroup.DELETE("/:id", handlers.DeleteUser)
	userGroup.GET("", handlers.ListUsers)
	userGroup.GET("/:id", handlers.GetUser)
}
