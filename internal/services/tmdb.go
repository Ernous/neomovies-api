package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"neomovies-api/internal/models"
)

type TMDBService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewTMDBService(apiKey string) *TMDBService {
	return &TMDBService{
		apiKey:  apiKey,
		baseURL: "https://api.themoviedb.org/3",
		client:  &http.Client{},
	}
}

func (s *TMDBService) SearchMovies(query string, page int, language, region string, year int) (*models.TMDBResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("query", query)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}
	
	if region != "" {
		params.Set("region", region)
	}
	
	if year > 0 {
		params.Set("year", strconv.Itoa(year))
	}

	endpoint := fmt.Sprintf("%s/search/movie?%s", s.baseURL, params.Encode())
	
	var response models.TMDBResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) SearchTVShows(query string, page int, language string, firstAirDateYear int) (*models.TMDBTVResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("query", query)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}
	
	if firstAirDateYear > 0 {
		params.Set("first_air_date_year", strconv.Itoa(firstAirDateYear))
	}

	endpoint := fmt.Sprintf("%s/search/tv?%s", s.baseURL, params.Encode())
	
	var response models.TMDBTVResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetMovie(id int, language string) (*models.Movie, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/movie/%d?%s", s.baseURL, id, params.Encode())
	
	var movie models.Movie
	err := s.makeRequest(endpoint, &movie)
	return &movie, err
}

func (s *TMDBService) GetTVShow(id int, language string) (*models.TVShow, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/tv/%d?%s", s.baseURL, id, params.Encode())
	
	var tvShow models.TVShow
	err := s.makeRequest(endpoint, &tvShow)
	return &tvShow, err
}

func (s *TMDBService) GetPopularMovies(page int, language, region string) (*models.TMDBResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}
	
	if region != "" {
		params.Set("region", region)
	}

	endpoint := fmt.Sprintf("%s/movie/popular?%s", s.baseURL, params.Encode())
	
	var response models.TMDBResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetTopRatedMovies(page int, language, region string) (*models.TMDBResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}
	
	if region != "" {
		params.Set("region", region)
	}

	endpoint := fmt.Sprintf("%s/movie/top_rated?%s", s.baseURL, params.Encode())
	
	var response models.TMDBResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetUpcomingMovies(page int, language, region string) (*models.TMDBResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}
	
	if region != "" {
		params.Set("region", region)
	}

	endpoint := fmt.Sprintf("%s/movie/upcoming?%s", s.baseURL, params.Encode())
	
	var response models.TMDBResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetNowPlayingMovies(page int, language, region string) (*models.TMDBResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}
	
	if region != "" {
		params.Set("region", region)
	}

	endpoint := fmt.Sprintf("%s/movie/now_playing?%s", s.baseURL, params.Encode())
	
	var response models.TMDBResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetMovieRecommendations(id, page int, language string) (*models.TMDBResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/movie/%d/recommendations?%s", s.baseURL, id, params.Encode())
	
	var response models.TMDBResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetSimilarMovies(id, page int, language string) (*models.TMDBResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/movie/%d/similar?%s", s.baseURL, id, params.Encode())
	
	var response models.TMDBResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetPopularTVShows(page int, language string) (*models.TMDBTVResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/tv/popular?%s", s.baseURL, params.Encode())
	
	var response models.TMDBTVResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetTopRatedTVShows(page int, language string) (*models.TMDBTVResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/tv/top_rated?%s", s.baseURL, params.Encode())
	
	var response models.TMDBTVResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetOnTheAirTVShows(page int, language string) (*models.TMDBTVResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/tv/on_the_air?%s", s.baseURL, params.Encode())
	
	var response models.TMDBTVResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetAiringTodayTVShows(page int, language string) (*models.TMDBTVResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/tv/airing_today?%s", s.baseURL, params.Encode())
	
	var response models.TMDBTVResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetTVRecommendations(id, page int, language string) (*models.TMDBTVResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/tv/%d/recommendations?%s", s.baseURL, id, params.Encode())
	
	var response models.TMDBTVResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) GetSimilarTVShows(id, page int, language string) (*models.TMDBTVResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))
	
	if language != "" {
		params.Set("language", language)
	} else {
		params.Set("language", "ru-RU")
	}

	endpoint := fmt.Sprintf("%s/tv/%d/similar?%s", s.baseURL, id, params.Encode())
	
	var response models.TMDBTVResponse
	err := s.makeRequest(endpoint, &response)
	return &response, err
}

func (s *TMDBService) makeRequest(endpoint string, target interface{}) error {
	resp, err := s.client.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}