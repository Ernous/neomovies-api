package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"neomovies-api/pkg/models"
)

// HealthCheck проверяет состояние API
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

	response := models.APIResponse{
		Success: true,
		Message: "API is running",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"environment_debug": envDebug,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}