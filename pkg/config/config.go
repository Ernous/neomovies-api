package config

import (
	"log"
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
	// Добавляем отладочное логирование для Vercel
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017/neomovies")
	log.Printf("DEBUG: MONGO_URI = %s (length: %d)", mongoURI, len(mongoURI))
	
	// Попробуем также другие возможные названия переменных
	if mongoURI == "mongodb://localhost:27017/neomovies" {
		alternatives := []string{"MONGODB_URI", "DATABASE_URL", "MONGO_URL"}
		for _, alt := range alternatives {
			if altValue := os.Getenv(alt); altValue != "" {
				log.Printf("DEBUG: Found alternative env var %s = %s", alt, altValue)
				mongoURI = altValue
				break
			}
		}
	}
	
	return &Config{
		MongoURI:        mongoURI,
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