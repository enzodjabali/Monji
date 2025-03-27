package routes

import (
	"monji/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Public routes.
	RegisterAuthRoutes(router)

	// Protected routes.
	api := router.Group("/")
	// Individual route groups add their own middleware as needed.

	RegisterEnvironmentRoutes(api)
	RegisterDatabaseRoutes(api)
	RegisterCollectionRoutes(api)
	RegisterDocumentRoutes(api)
	RegisterMongoUserRoutes(api)
	RegisterWhoAmIRoute(api)
	RegisterUserRoutes(api)
	RegisterPermissionsRoutes(api)

	return router
}
