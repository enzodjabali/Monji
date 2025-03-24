package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterEnvironmentRoutes(rg *gin.RouterGroup) {
	envGroup := rg.Group("/environments")
	{
		envGroup.GET("", middleware.AdminMiddleware(), handlers.ListEnvironments)
		envGroup.GET("/:id", middleware.AdminMiddleware(), handlers.GetEnvironment)
		envGroup.POST("", middleware.AdminMiddleware(), handlers.CreateEnvironment)
		envGroup.PUT("/:id", middleware.AdminMiddleware(), handlers.UpdateEnvironment)
		envGroup.DELETE("/:id", middleware.AdminMiddleware(), handlers.DeleteEnvironment)
	}
}
