package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterEnvironmentRoutes requires only authentication, not admin role.
// The in-handler logic checks if a user has the correct environment permission.
func RegisterEnvironmentRoutes(rg *gin.RouterGroup) {
	envGroup := rg.Group("/environments")
	// Only require authentication
	envGroup.Use(middleware.AuthMiddleware())

	// Then let the handlers check read/write access
	envGroup.GET("", handlers.ListEnvironments)
	envGroup.GET("/:id", handlers.GetEnvironment)
	envGroup.POST("", handlers.CreateEnvironment)
	envGroup.PUT("/:id", handlers.UpdateEnvironment)
	envGroup.DELETE("/:id", handlers.DeleteEnvironment)
}
