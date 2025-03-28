package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, proceeding with system environment variables")
		}
	}

	cfg := &Config{
		Port:       os.Getenv("PORT"),
		SQLitePath: os.Getenv("SQLITE_PATH"),
		MongoURI:   os.Getenv("MONGO_URI"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
	if cfg.Port == "" {
		return nil, errors.New("environment variable PORT is not set")
	}
	if cfg.SQLitePath == "" {
		return nil, errors.New("environment variable SQLITE_PATH is not set")
	}
	if cfg.MongoURI == "" {
		return nil, errors.New("environment variable MONGO_URI is not set")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("environment variable JWT_SECRET is not set")
	}
	return cfg, nil
}
