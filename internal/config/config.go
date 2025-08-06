package config

import (
	"os"
)

type Config struct {
	MongoURI   string
	TMDBAPIKey string
	JWTSecret  string
	Port       string
	BaseURL    string
}

func New() *Config {
	return &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017/neomovies"),
		TMDBAPIKey: getEnv("TMDB_API_KEY", ""),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		Port:       getEnv("PORT", "3000"),
		BaseURL:    getEnv("BASE_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}