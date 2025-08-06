package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"neomovies-api/internal/models"
)

type TorrentService struct {
	client *http.Client
}

func NewTorrentService() *TorrentService {
	return &TorrentService{
		client: &http.Client{},
	}
}

type TorrentSearchOptions struct {
	Season       *int
	Quality      []string
	MinQuality   string
	MaxQuality   string
	ExcludeQualities []string
	HDR          *bool
	HEVC         *bool
	SortBy       string
	SortOrder    string
	GroupByQuality bool
	GroupBySeason  bool
}

func (s *TorrentService) SearchTorrentsByIMDbID(tmdbService *TMDBService, imdbID, mediaType string, options *TorrentSearchOptions) (*models.TorrentSearchResponse, error) {
	// Получаем информацию о фильме/сериале из TMDB
	title, year, err := s.getTitleFromTMDB(tmdbService, imdbID, mediaType)
	if err != nil {
		return nil, fmt.Errorf("failed to get title from TMDB: %w", err)
	}

	// Поиск торрентов на bitru.org
	torrents, err := s.searchOnBitru(title, year, mediaType, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search torrents: %w", err)
	}

	// Фильтрация и сортировка
	filteredTorrents := s.filterTorrents(torrents, options)
	sortedTorrents := s.sortTorrents(filteredTorrents, options.SortBy, options.SortOrder)

	return &models.TorrentSearchResponse{
		Query:   fmt.Sprintf("%s (%s)", title, year),
		Results: sortedTorrents,
		Total:   len(sortedTorrents),
	}, nil
}

func (s *TorrentService) getTitleFromTMDB(tmdbService *TMDBService, imdbID, mediaType string) (string, string, error) {
	// Используем find API для поиска по IMDB ID
	endpoint := fmt.Sprintf("https://api.themoviedb.org/3/find/%s", imdbID)
	
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", "", err
	}

	// Добавляем параметры
	params := url.Values{}
	params.Set("external_source", "imdb_id")
	params.Set("language", "ru-RU")
	req.URL.RawQuery = params.Encode()

	// Добавляем авторизацию
	req.Header.Set("Authorization", "Bearer "+tmdbService.accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var findResponse struct {
		MovieResults []struct {
			Title         string `json:"title"`
			OriginalTitle string `json:"original_title"`
			ReleaseDate   string `json:"release_date"`
		} `json:"movie_results"`
		TVResults []struct {
			Name         string `json:"name"`
			OriginalName string `json:"original_name"`
			FirstAirDate string `json:"first_air_date"`
		} `json:"tv_results"`
	}

	if err := json.Unmarshal(body, &findResponse); err != nil {
		return "", "", err
	}

	if mediaType == "movie" && len(findResponse.MovieResults) > 0 {
		movie := findResponse.MovieResults[0]
		title := movie.OriginalTitle
		if title == "" {
			title = movie.Title
		}
		year := ""
		if movie.ReleaseDate != "" {
			year = movie.ReleaseDate[:4]
		}
		return title, year, nil
	}

	if (mediaType == "tv" || mediaType == "series") && len(findResponse.TVResults) > 0 {
		tv := findResponse.TVResults[0]
		title := tv.OriginalName
		if title == "" {
			title = tv.Name
		}
		year := ""
		if tv.FirstAirDate != "" {
			year = tv.FirstAirDate[:4]
		}
		return title, year, nil
	}

	return "", "", fmt.Errorf("no results found for IMDB ID: %s", imdbID)
}

func (s *TorrentService) searchOnBitru(title, year, mediaType string, options *TorrentSearchOptions) ([]models.TorrentResult, error) {
	// Формируем поисковый запрос
	searchQuery := title
	if year != "" {
		searchQuery += " " + year
	}

	// URL для поиска на bitru.org
	searchURL := fmt.Sprintf("https://bitru.org/search.php?search=%s", url.QueryEscape(searchQuery))

	resp, err := s.client.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Парсим HTML и извлекаем торренты
	torrents := s.parseHTML(string(body))
	
	return torrents, nil
}

func (s *TorrentService) parseHTML(html string) []models.TorrentResult {
	var torrents []models.TorrentResult

	// Простой парсинг HTML (в реальном проекте лучше использовать goquery)
	// Ищем строки таблицы с торрентами
	re := regexp.MustCompile(`<tr[^>]*>.*?</tr>`)
	rows := re.FindAllString(html, -1)

	for _, row := range rows {
		torrent := s.parseTorrentRow(row)
		if torrent.Title != "" {
			torrents = append(torrents, torrent)
		}
	}

	return torrents
}

func (s *TorrentService) parseTorrentRow(row string) models.TorrentResult {
	var torrent models.TorrentResult

	// Извлекаем название
	titleRe := regexp.MustCompile(`<a[^>]*href="[^"]*"[^>]*>([^<]+)</a>`)
	if matches := titleRe.FindStringSubmatch(row); len(matches) > 1 {
		torrent.Title = strings.TrimSpace(matches[1])
	}

	// Извлекаем размер
	sizeRe := regexp.MustCompile(`(\d+(?:\.\d+)?\s*[KMGTPE]?B)`)
	if matches := sizeRe.FindStringSubmatch(row); len(matches) > 1 {
		torrent.Size = matches[1]
	}

	// Извлекаем количество сидеров и личеров
	seedRe := regexp.MustCompile(`class="seed[^"]*"[^>]*>(\d+)</`)
	if matches := seedRe.FindStringSubmatch(row); len(matches) > 1 {
		if seeders, err := strconv.Atoi(matches[1]); err == nil {
			torrent.Seeders = seeders
		}
	}

	leechRe := regexp.MustCompile(`class="leech[^"]*"[^>]*>(\d+)</`)
	if matches := leechRe.FindStringSubmatch(row); len(matches) > 1 {
		if leechers, err := strconv.Atoi(matches[1]); err == nil {
			torrent.Leechers = leechers
		}
	}

	// Определяем качество по названию
	torrent.Quality = s.extractQuality(torrent.Title)

	// Извлекаем magnet ссылку
	magnetRe := regexp.MustCompile(`href="(magnet:[^"]+)"`)
	if matches := magnetRe.FindStringSubmatch(row); len(matches) > 1 {
		torrent.MagnetLink = matches[1]
	}

	return torrent
}

func (s *TorrentService) extractQuality(title string) string {
	title = strings.ToUpper(title)
	
	qualities := []string{"2160P", "4K", "1440P", "1080P", "720P", "480P", "360P"}
	for _, quality := range qualities {
		if strings.Contains(title, quality) {
			if quality == "2160P" {
				return "4K"
			}
			return quality
		}
	}
	
	return "Unknown"
}

func (s *TorrentService) filterTorrents(torrents []models.TorrentResult, options *TorrentSearchOptions) []models.TorrentResult {
	if options == nil {
		return torrents
	}

	var filtered []models.TorrentResult

	for _, torrent := range torrents {
		// Фильтрация по качеству
		if len(options.Quality) > 0 {
			found := false
			for _, quality := range options.Quality {
				if strings.EqualFold(torrent.Quality, quality) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Фильтрация по минимальному качеству
		if options.MinQuality != "" && !s.qualityMeetsMinimum(torrent.Quality, options.MinQuality) {
			continue
		}

		// Фильтрация по максимальному качеству
		if options.MaxQuality != "" && !s.qualityMeetsMaximum(torrent.Quality, options.MaxQuality) {
			continue
		}

		// Исключение качеств
		if len(options.ExcludeQualities) > 0 {
			excluded := false
			for _, excludeQuality := range options.ExcludeQualities {
				if strings.EqualFold(torrent.Quality, excludeQuality) {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
		}

		// Фильтрация по HDR
		if options.HDR != nil {
			hasHDR := strings.Contains(strings.ToUpper(torrent.Title), "HDR")
			if *options.HDR != hasHDR {
				continue
			}
		}

		// Фильтрация по HEVC
		if options.HEVC != nil {
			hasHEVC := strings.Contains(strings.ToUpper(torrent.Title), "HEVC") || 
					  strings.Contains(strings.ToUpper(torrent.Title), "H.265")
			if *options.HEVC != hasHEVC {
				continue
			}
		}

		filtered = append(filtered, torrent)
	}

	return filtered
}

func (s *TorrentService) qualityMeetsMinimum(quality, minQuality string) bool {
	qualityOrder := map[string]int{
		"360P": 1, "480P": 2, "720P": 3, "1080P": 4, "1440P": 5, "4K": 6, "2160P": 6,
	}
	
	currentLevel := qualityOrder[strings.ToUpper(quality)]
	minLevel := qualityOrder[strings.ToUpper(minQuality)]
	
	return currentLevel >= minLevel
}

func (s *TorrentService) qualityMeetsMaximum(quality, maxQuality string) bool {
	qualityOrder := map[string]int{
		"360P": 1, "480P": 2, "720P": 3, "1080P": 4, "1440P": 5, "4K": 6, "2160P": 6,
	}
	
	currentLevel := qualityOrder[strings.ToUpper(quality)]
	maxLevel := qualityOrder[strings.ToUpper(maxQuality)]
	
	return currentLevel <= maxLevel
}

func (s *TorrentService) sortTorrents(torrents []models.TorrentResult, sortBy, sortOrder string) []models.TorrentResult {
	if sortBy == "" {
		sortBy = "seeders"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	sort.Slice(torrents, func(i, j int) bool {
		var less bool
		
		switch sortBy {
		case "seeders":
			less = torrents[i].Seeders < torrents[j].Seeders
		case "size":
			less = s.compareSizes(torrents[i].Size, torrents[j].Size)
		case "date":
			less = torrents[i].AddedDate < torrents[j].AddedDate
		default:
			less = torrents[i].Seeders < torrents[j].Seeders
		}

		if sortOrder == "asc" {
			return less
		}
		return !less
	})

	return torrents
}

func (s *TorrentService) compareSizes(size1, size2 string) bool {
	// Простое сравнение размеров (можно улучшить)
	return len(size1) < len(size2)
}