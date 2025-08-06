package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"neomovies-api/internal/models"
	"neomovies-api/internal/services"
)

type TorrentsHandler struct {
	torrentService *services.TorrentService
	tmdbService    *services.TMDBService
}

func NewTorrentsHandler(torrentService *services.TorrentService, tmdbService *services.TMDBService) *TorrentsHandler {
	return &TorrentsHandler{
		torrentService: torrentService,
		tmdbService:    tmdbService,
	}
}

func (h *TorrentsHandler) SearchTorrents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imdbID := vars["imdbId"]

	if imdbID == "" {
		http.Error(w, "IMDB ID is required", http.StatusBadRequest)
		return
	}

	// Параметры запроса
	mediaType := r.URL.Query().Get("type")
	if mediaType == "" {
		mediaType = "movie"
	}

	// Создаем опции поиска
	options := &services.TorrentSearchOptions{}

	// Качество
	if quality := r.URL.Query().Get("quality"); quality != "" {
		options.Quality = strings.Split(quality, ",")
	}

	// Минимальное и максимальное качество
	options.MinQuality = r.URL.Query().Get("minQuality")
	options.MaxQuality = r.URL.Query().Get("maxQuality")

	// Исключаемые качества
	if excludeQualities := r.URL.Query().Get("excludeQualities"); excludeQualities != "" {
		options.ExcludeQualities = strings.Split(excludeQualities, ",")
	}

	// HDR
	if hdr := r.URL.Query().Get("hdr"); hdr != "" {
		if hdrBool, err := strconv.ParseBool(hdr); err == nil {
			options.HDR = &hdrBool
		}
	}

	// HEVC
	if hevc := r.URL.Query().Get("hevc"); hevc != "" {
		if hevcBool, err := strconv.ParseBool(hevc); err == nil {
			options.HEVC = &hevcBool
		}
	}

	// Сортировка
	options.SortBy = r.URL.Query().Get("sortBy")
	if options.SortBy == "" {
		options.SortBy = "seeders"
	}

	options.SortOrder = r.URL.Query().Get("sortOrder")
	if options.SortOrder == "" {
		options.SortOrder = "desc"
	}

	// Группировка
	if groupByQuality := r.URL.Query().Get("groupByQuality"); groupByQuality == "true" {
		options.GroupByQuality = true
	}

	if groupBySeason := r.URL.Query().Get("groupBySeason"); groupBySeason == "true" {
		options.GroupBySeason = true
	}

	// Сезон для сериалов
	if season := r.URL.Query().Get("season"); season != "" {
		if seasonInt, err := strconv.Atoi(season); err == nil {
			options.Season = &seasonInt
		}
	}

	// Поиск торрентов
	results, err := h.torrentService.SearchTorrentsByIMDbID(h.tmdbService, imdbID, mediaType, options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(results.Results) == 0 {
		response := map[string]interface{}{
			"imdbId": imdbID,
			"type":   mediaType,
			"total":  0,
			"results": []models.TorrentResult{},
			"error":  "No torrents found for this IMDB ID",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Формируем ответ
	response := map[string]interface{}{
		"imdbId":  imdbID,
		"type":    mediaType,
		"total":   results.Total,
		"grouped": options.GroupByQuality || options.GroupBySeason,
		"results": results.Results,
	}

	if options.Season != nil {
		response["season"] = *options.Season
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    response,
	})
}