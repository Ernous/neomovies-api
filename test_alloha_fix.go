package main

import (
	"fmt"
	"log"

	"neomovies-api/pkg/services"
)

func main() {
	// Создаем торрент-сервис
	torrentService := services.NewTorrentService()
	
	// Создаем пустой TMDB сервис (не будет использоваться благодаря fallback на Alloha)
	tmdbService := services.NewTMDBService("")
	
	// Тестируем проблемный IMDB ID tt6385540 (сериал Hilda)
	imdbID := "tt6385540"
	mediaType := "tv"
	
	fmt.Printf("Тестируем IMDB ID: %s (тип: %s)\n", imdbID, mediaType)
	
	// Тестируем поиск торрентов
	fmt.Printf("Тестируем поиск торрентов...\n")
	options := &services.TorrentSearchOptions{}
	response, err := torrentService.SearchTorrentsByIMDbID(tmdbService, imdbID, mediaType, options)
	if err != nil {
		log.Printf("Ошибка поиска торрентов: %v", err)
		return
	}
	
	fmt.Printf("Найдено торрентов: %d\n", response.Total)
	for i, torrent := range response.Results {
		if i >= 3 { // Показываем только первые 3
			break
		}
		fmt.Printf("  %d. %s (Seeders: %d)\n", i+1, torrent.Title, torrent.Seeders)
	}
}