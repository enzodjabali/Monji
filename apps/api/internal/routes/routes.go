package routes

import (
	"monji/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Public routes.
	RegisterAuthRoutes(router) // /login, etc.

	// Protected routes group:
	api := router.Group("/")

	// We register each route set. Each set includes AuthMiddleware in its own file if needed.
	RegisterEnvironmentRoutes(api)
	RegisterDatabaseRoutes(api)
	RegisterCollectionRoutes(api)
	RegisterDocumentRoutes(api)
	RegisterMongoUserRoutes(api)
	RegisterUserRoutes(api)        // userGroup still has AdminMiddleware
	RegisterPermissionsRoutes(api) // presumably also admin only
	RegisterWhoAmIRoute(api)

	return router
}
