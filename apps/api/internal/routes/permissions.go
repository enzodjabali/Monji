package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPermissionsRoutes allows an admin or superadmin to set environment/db permissions for normal users.
//
// POST /users/:userId/environments/:envId/permissions
// POST /users/:userId/environments/:envId/databases/:dbName/permissions
func RegisterPermissionsRoutes(rg *gin.RouterGroup) {
	// Only admin or superadmin can change user permissions
	adminGroup := rg.Group("/users/:userId")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

	// environment-level
	adminGroup.POST("/environments/:envId/permissions", handlers.SetUserEnvironmentPermission)
	// database-level
	adminGroup.POST("/environments/:envId/databases/:dbName/permissions", handlers.SetUserDBPermission)
}
