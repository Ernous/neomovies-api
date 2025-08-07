package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"neomovies-api/pkg/config"
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
	var allohaResponse struct {
		Status string `json:"status"`
		Data struct {
			Iframe string `json:"iframe"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &allohaResponse); err != nil {
		http.Error(w, "Invalid JSON from Alloha", http.StatusBadGateway)
		return
	}
	if allohaResponse.Status != "success" || allohaResponse.Data.Iframe == "" {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	iframeCode := allohaResponse.Data.Iframe
	if !strings.Contains(iframeCode, "<") {
		iframeCode = fmt.Sprintf(`<iframe src="%s" allowfullscreen style="border:none;width:100%%;height:100%%"></iframe>`, iframeCode)
	}
	htmlDoc := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset='utf-8'/><title>Alloha Player</title><style>html,body{margin:0;height:100%%;}</style></head><body>%s</body></html>`, iframeCode)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlDoc))
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
	url := fmt.Sprintf("%s?imdb_id=%s", h.config.LumexURL, url.QueryEscape(imdbID))
	iframe := fmt.Sprintf(`<iframe src="%s" allowfullscreen loading="lazy" style="border:none;width:100%%;height:100%%;"></iframe>`, url)
	htmlDoc := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset='utf-8'/><title>Lumex Player</title><style>html,body{margin:0;height:100%%;}</style></head><body>%s</body></html>`, iframe)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlDoc))
}