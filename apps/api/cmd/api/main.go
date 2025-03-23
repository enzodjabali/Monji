package main

import (
	"context"
	"log"
	"time"

	"monji/internal/config"
	"monji/internal/database"
	"monji/internal/router"
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

	// Set up the router.
	r := router.SetupRouter(cfg)

	// Start the server.
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
