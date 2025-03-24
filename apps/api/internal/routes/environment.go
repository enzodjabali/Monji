package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterEnvironmentRoutes requires Auth + Admin for all environment endpoints.
func RegisterEnvironmentRoutes(rg *gin.RouterGroup) {
	envGroup := rg.Group("/environments")
	{
		// Apply both AuthMiddleware and AdminMiddleware to the group:
		envGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

		envGroup.GET("", handlers.ListEnvironments)
		envGroup.GET("/:id", handlers.GetEnvironment)
		envGroup.POST("", handlers.CreateEnvironment)
		envGroup.PUT("/:id", handlers.UpdateEnvironment)
		envGroup.DELETE("/:id", handlers.DeleteEnvironment)
	}
}
