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
    htmlDoc := `<!DOCTYPE html><html><head><meta charset='utf-8'/><title>Test</title></head><body><iframe src="https://www.youtube.com/embed/dQw4w9WgXcQ" allowfullscreen style="border:none;width:100%;height:100%"></iframe></body></html>`
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