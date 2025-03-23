package main

import (
	"context"
	"log"
	"time"

	"monji/internal/config"
	"monji/internal/database"
	"monji/internal/handlers"
	"monji/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize SQLite (for users and environments).
	database.InitSQLite(cfg.SQLitePath)

	// (Optional) Initialize a global MongoDB connection if needed.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = database.ConnectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Printf("Warning: Failed to connect to main MongoDB instance: %v", err)
	}

	// Create the Gin router.
	router := gin.Default()

	// Public route.
	router.POST("/login", handlers.Login)

	// Protected routes.
	api := router.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		// Environment endpoints.
		api.GET("/environments", handlers.ListEnvironments)
		api.POST("/environments", middleware.AdminMiddleware(), handlers.CreateEnvironment)
		api.PUT("/environments/:id", middleware.AdminMiddleware(), handlers.UpdateEnvironment)
		api.DELETE("/environments/:id", middleware.AdminMiddleware(), handlers.DeleteEnvironment)

		// Mongo queries.
		api.GET("/environments/:id/databases", handlers.GetDatabases)
		api.GET("/environments/:id/databases/:dbName/collections", handlers.GetCollections)
		api.GET("/environments/:id/databases/:dbName/collections/:collName/documents", handlers.GetDocuments)
	}

	// Start the server.
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
