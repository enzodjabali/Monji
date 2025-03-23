package config

import "os"

// Config holds configuration values.
type Config struct {
	Port       string
	SQLitePath string
	MongoURI   string
	JWTSecret  string
}

// LoadConfig reads config from environment variables with defaults.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		SQLitePath: getEnv("SQLITE_PATH", "/data/sqlite/users.db"),
		MongoURI:   getEnv("MONGO_URI", "mongodb://root:strongpassword@mongodb:27017"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
