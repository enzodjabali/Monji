package routes

import (
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterWhoAmIRoute sets up the /whoami endpoint to return the current user details.
func RegisterWhoAmIRoute(rg *gin.RouterGroup) {
	// Use AuthMiddleware to ensure the user is authenticated.
	rg.GET("/whoami", middleware.AuthMiddleware(), handlers.WhoAmI)
}
