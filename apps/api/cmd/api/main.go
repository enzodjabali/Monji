package main

import (
	"context"
	"log"
	"time"

	"monji/internal/config"
	"monji/internal/database"
	"monji/internal/routes"
)

func main() {
	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize SQLite (for users and environments).
	database.InitSQLite(cfg.SQLitePath)

	// (Optional) Test MongoDB connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = database.ConnectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Printf("Warning: Failed to connect to MongoDB: %v", err)
	}

	// Set up all routes.
	router := routes.SetupRoutes(cfg)

	// Start the server.
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
