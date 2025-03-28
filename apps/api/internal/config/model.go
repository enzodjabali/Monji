package config

type Config struct {
	Port       string
	SQLitePath string
	MongoURI   string
	JWTSecret  string
}
