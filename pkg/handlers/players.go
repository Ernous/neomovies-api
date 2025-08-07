package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"neomovies-api/pkg/config"
	"neomovies-api/pkg/models"
	"github.com/gorilla/mux"
)

type PlayersHandler struct {
	config *config.Config
}

func NewPlayersHandler(cfg *config.Config) *PlayersHandler {
	return &PlayersHandler{
		config: cfg,
	}
}

func (h *PlayersHandler) GetAllohaPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imdbID := vars["imdb_id"]
	if imdbID == "" {
		http.Error(w, "imdb_id path param is required", http.StatusBadRequest)
		return
	}
	if h.config.AllohaToken == "" {
		http.Error(w, "Server misconfiguration: ALLOHA_TOKEN missing", http.StatusInternalServerError)
		return
	}
	idParam := fmt.Sprintf("imdb=%s", url.QueryEscape(imdbID))
	apiURL := fmt.Sprintf("https://api.alloha.tv/?token=%s&%s", h.config.AllohaToken, idParam)
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch from Alloha API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Alloha API error: %d", resp.StatusCode), http.StatusBadGateway)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Alloha response", http.StatusInternalServerError)
		return
	}
	var allohaResponse map[string]interface{}
	if err := json.Unmarshal(body, &allohaResponse); err != nil {
		http.Error(w, "Invalid JSON from Alloha", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    allohaResponse,
	})
}

func (h *PlayersHandler) GetLumexPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imdbID := vars["imdb_id"]
	if imdbID == "" {
		http.Error(w, "imdb_id path param is required", http.StatusBadRequest)
		return
	}
	if h.config.LumexURL == "" {
		http.Error(w, "Server misconfiguration: LUMEX_URL missing", http.StatusInternalServerError)
		return
	}
	params := url.Values{}
	params.Set("imdb_id", imdbID)
	apiURL := fmt.Sprintf("%s?%s", h.config.LumexURL, params.Encode())
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch from Lumex API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Lumex API error: %d", resp.StatusCode), http.StatusBadGateway)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Lumex response", http.StatusInternalServerError)
		return
	}
	var lumexResponse map[string]interface{}
	if err := json.Unmarshal(body, &lumexResponse); err != nil {
		http.Error(w, "Invalid JSON from Lumex", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    lumexResponse,
	})
}