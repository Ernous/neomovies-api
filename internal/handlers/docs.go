package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type DocsHandler struct {
	openAPISpec *OpenAPISpec
}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{
		openAPISpec: getOpenAPISpec(),
	}
}

func (h *DocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Обслуживаем статические файлы для документации
	if r.URL.Path == "/" {
		h.serveDocs(w, r)
		return
	}
	
	http.NotFound(w, r)
}

func (h *DocsHandler) RedirectToDocs(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
}

func (h *DocsHandler) GetOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(h.openAPISpec)
}

func (h *DocsHandler) serveDocs(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		if r.TLS != nil {
			baseURL = fmt.Sprintf("https://%s", r.Host)
		} else {
			baseURL = fmt.Sprintf("http://%s", r.Host)
		}
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Neo Movies API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <script
        id="api-reference"
        data-url="{{.BaseURL}}/openapi.json"
        data-configuration='{
            "theme": "saturn",
            "layout": "modern",
            "defaultHttpClient": {
                "targetKey": "javascript",
                "clientKey": "fetch"
            },
            "authentication": {
                "securitySchemes": {
                    "bearerAuth": {
                        "type": "http",
                        "scheme": "bearer",
                        "bearerFormat": "JWT"
                    }
                }
            },
            "spec": {
                "url": "{{.BaseURL}}/openapi.json"
            },
            "metadata": {
                "title": "Neo Movies API",
                "description": "Современный API для поиска фильмов и сериалов с поддержкой авторизации",
                "favicon": "https://cdn.jsdelivr.net/npm/@scalar/api-reference/dist/browser/favicon.ico"
            }
        }'></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@latest"></script>
</body>
</html>`

	t, err := template.New("docs").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, map[string]string{
		"BaseURL": baseURL,
	})
}

type OpenAPISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    Info                   `json:"info"`
	Servers []Server               `json:"servers"`
	Paths   map[string]interface{} `json:"paths"`
	Components Components           `json:"components"`
}

type Info struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	Contact     Contact `json:"contact"`
}

type Contact struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type Components struct {
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
	Schemas         map[string]interface{}    `json:"schemas"`
}

type SecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

func getOpenAPISpec() *OpenAPISpec {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	return &OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "Neo Movies API",
			Description: "Современный API для поиска фильмов и сериалов с интеграцией TMDB и поддержкой авторизации",
			Version:     "2.0.0",
			Contact: Contact{
				Name: "API Support",
				URL:  "https://github.com/your-username/neomovies-api-go",
			},
		},
		Servers: []Server{
			{
				URL:         baseURL,
				Description: "Production server",
			},
		},
		Paths: map[string]interface{}{
			"/api/v1/health": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Health Check",
					"description": "Проверка работоспособности API",
					"tags": []string{"Health"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "API работает корректно",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/APIResponse",
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/auth/register": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Регистрация пользователя",
					"description": "Создание нового аккаунта пользователя",
					"tags": []string{"Authentication"},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/RegisterRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Пользователь успешно зарегистрирован",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/AuthResponse",
									},
								},
							},
						},
						"409": map[string]interface{}{
							"description": "Пользователь с таким email уже существует",
						},
					},
				},
			},
			"/api/v1/auth/login": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Авторизация пользователя",
					"description": "Получение JWT токена для доступа к приватным эндпоинтам",
					"tags": []string{"Authentication"},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/LoginRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Успешная авторизация",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/AuthResponse",
									},
								},
							},
						},
						"401": map[string]interface{}{
							"description": "Неверный email или пароль",
						},
					},
				},
			},
			"/api/v1/auth/profile": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Получить профиль пользователя",
					"description": "Получение информации о текущем пользователе",
					"tags": []string{"Authentication"},
					"security": []map[string][]string{
						{"bearerAuth": []string{}},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Информация о пользователе",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/User",
									},
								},
							},
						},
					},
				},
				"put": map[string]interface{}{
					"summary": "Обновить профиль пользователя",
					"description": "Обновление информации о пользователе",
					"tags": []string{"Authentication"},
					"security": []map[string][]string{
						{"bearerAuth": []string{}},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Профиль успешно обновлен",
						},
					},
				},
			},
			"/api/v1/movies/search": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Поиск фильмов",
					"description": "Поиск фильмов по названию с поддержкой фильтров",
					"tags": []string{"Movies"},
					"parameters": []map[string]interface{}{
						{
							"name": "query",
							"in": "query",
							"required": true,
							"schema": map[string]string{"type": "string"},
							"description": "Поисковый запрос",
						},
						{
							"name": "page",
							"in": "query",
							"schema": map[string]string{"type": "integer", "default": "1"},
							"description": "Номер страницы",
						},
						{
							"name": "language",
							"in": "query",
							"schema": map[string]string{"type": "string", "default": "ru-RU"},
							"description": "Язык ответа",
						},
						{
							"name": "year",
							"in": "query",
							"schema": map[string]string{"type": "integer"},
							"description": "Год выпуска",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Результаты поиска фильмов",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/MovieSearchResponse",
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/movies/popular": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Популярные фильмы",
					"description": "Получение списка популярных фильмов",
					"tags": []string{"Movies"},
					"parameters": []map[string]interface{}{
						{
							"name": "page",
							"in": "query",
							"schema": map[string]string{"type": "integer", "default": "1"},
						},
						{
							"name": "language",
							"in": "query",
							"schema": map[string]string{"type": "string", "default": "ru-RU"},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Список популярных фильмов",
						},
					},
				},
			},
			"/api/v1/movies/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Получить фильм по ID",
					"description": "Подробная информация о фильме",
					"tags": []string{"Movies"},
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]string{"type": "integer"},
							"description": "ID фильма в TMDB",
						},
						{
							"name": "language",
							"in": "query",
							"schema": map[string]string{"type": "string", "default": "ru-RU"},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Информация о фильме",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/Movie",
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/favorites": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Получить избранные фильмы",
					"description": "Список избранных фильмов пользователя",
					"tags": []string{"Favorites"},
					"security": []map[string][]string{
						{"bearerAuth": []string{}},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Список избранных фильмов",
						},
					},
				},
			},
			"/api/v1/favorites/{id}": map[string]interface{}{
				"post": map[string]interface{}{
					"summary": "Добавить в избранное",
					"description": "Добавление фильма в избранное",
					"tags": []string{"Favorites"},
					"security": []map[string][]string{
						{"bearerAuth": []string{}},
					},
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]string{"type": "string"},
							"description": "ID фильма",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Фильм добавлен в избранное",
						},
					},
				},
				"delete": map[string]interface{}{
					"summary": "Удалить из избранного",
					"description": "Удаление фильма из избранного",
					"tags": []string{"Favorites"},
					"security": []map[string][]string{
						{"bearerAuth": []string{}},
					},
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]string{"type": "string"},
							"description": "ID фильма",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Фильм удален из избранного",
						},
					},
				},
			},
		},
		Components: Components{
			SecuritySchemes: map[string]SecurityScheme{
				"bearerAuth": {
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
			Schemas: map[string]interface{}{
				"APIResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"success": map[string]string{"type": "boolean"},
						"data": map[string]string{"type": "object"},
						"message": map[string]string{"type": "string"},
						"error": map[string]string{"type": "string"},
					},
				},
				"RegisterRequest": map[string]interface{}{
					"type": "object",
					"required": []string{"email", "password", "name"},
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type": "string",
							"format": "email",
							"example": "user@example.com",
						},
						"password": map[string]interface{}{
							"type": "string",
							"minLength": 6,
							"example": "password123",
						},
						"name": map[string]interface{}{
							"type": "string",
							"example": "Иван Иванов",
						},
					},
				},
				"LoginRequest": map[string]interface{}{
					"type": "object",
					"required": []string{"email", "password"},
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type": "string",
							"format": "email",
							"example": "user@example.com",
						},
						"password": map[string]interface{}{
							"type": "string",
							"example": "password123",
						},
					},
				},
				"AuthResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"token": map[string]string{"type": "string"},
						"user": map[string]interface{}{"$ref": "#/components/schemas/User"},
					},
				},
				"User": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]string{"type": "string"},
						"email": map[string]string{"type": "string"},
						"name": map[string]string{"type": "string"},
						"avatar": map[string]string{"type": "string"},
						"favorites": map[string]interface{}{
							"type": "array",
							"items": map[string]string{"type": "string"},
						},
						"created_at": map[string]interface{}{
							"type": "string",
							"format": "date-time",
						},
						"updated_at": map[string]interface{}{
							"type": "string",
							"format": "date-time",
						},
					},
				},
				"Movie": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]string{"type": "integer"},
						"title": map[string]string{"type": "string"},
						"original_title": map[string]string{"type": "string"},
						"overview": map[string]string{"type": "string"},
						"poster_path": map[string]string{"type": "string"},
						"backdrop_path": map[string]string{"type": "string"},
						"release_date": map[string]string{"type": "string"},
						"vote_average": map[string]string{"type": "number"},
						"vote_count": map[string]string{"type": "integer"},
						"popularity": map[string]string{"type": "number"},
						"adult": map[string]string{"type": "boolean"},
						"original_language": map[string]string{"type": "string"},
					},
				},
				"MovieSearchResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"page": map[string]string{"type": "integer"},
						"results": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{"$ref": "#/components/schemas/Movie"},
						},
						"total_pages": map[string]string{"type": "integer"},
						"total_results": map[string]string{"type": "integer"},
					},
				},
			},
		},
	}
}