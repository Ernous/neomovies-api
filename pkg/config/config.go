package config

import (
	"os"
)

type Config struct {
	MongoURI        string
	TMDBAccessToken string
	JWTSecret       string
	Port            string
	BaseURL         string
	NodeEnv         string
	GmailUser       string
	GmailPassword   string
	LumexURL        string
	AllohaToken     string
}

func New() *Config {
	return &Config{
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017/neomovies"),
		TMDBAccessToken: getEnv("TMDB_ACCESS_TOKEN", ""),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		Port:            getEnv("PORT", "3000"),
		BaseURL:         getEnv("BASE_URL", "http://localhost:3000"),
		NodeEnv:         getEnv("NODE_ENV", "development"),
		GmailUser:       getEnv("GMAIL_USER", ""),
		GmailPassword:   getEnv("GMAIL_APP_PASSWORD", ""),
		LumexURL:        getEnv("LUMEX_URL", ""),
		AllohaToken:     getEnv("ALLOHA_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}