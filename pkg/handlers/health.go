package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"neomovies-api/pkg/models"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Добавляем отладочную информацию о переменных окружения
	envDebug := map[string]interface{}{
		"MONGO_URI_set":       os.Getenv("MONGO_URI") != "",
		"MONGODB_URI_set":     os.Getenv("MONGODB_URI") != "",
		"DATABASE_URL_set":    os.Getenv("DATABASE_URL") != "",
		"MONGO_URL_set":       os.Getenv("MONGO_URL") != "",
		"TMDB_ACCESS_TOKEN_set": os.Getenv("TMDB_ACCESS_TOKEN") != "",
		"JWT_SECRET_set":      os.Getenv("JWT_SECRET") != "",
		"PORT":               os.Getenv("PORT"),
		"NODE_ENV":           os.Getenv("NODE_ENV"),
	}

	health := map[string]interface{}{
		"status":    "OK",
		"timestamp": time.Now().UTC(),
		"service":   "neomovies-api",
		"version":   "2.0.0",
		"uptime":    time.Since(startTime),
		"environment_debug": envDebug,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Message: "API is running",
		Data:    health,
	})
}

var startTime = time.Now()