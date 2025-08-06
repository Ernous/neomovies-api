# Деплой на Vercel

## Исправленные проблемы

### ❌ Ошибка: "The 'functions' property cannot be used in conjunction with the 'builds' property"

**Проблема:** В `vercel.json` одновременно использовались свойства `builds` и `functions`, что не разрешено.

**Решение:** Удалено свойство `functions`, конфигурация `maxDuration` перенесена в `builds.config`.

### ✅ Исправленная конфигурация `vercel.json`:

```json
{
  "version": 2,
  "builds": [
    {
      "src": "api/index.go",
      "use": "@vercel/go",
      "config": {
        "maxDuration": 10
      }
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "/api/index.go"
    }
  ],
  "env": {
    "GO_VERSION": "1.21"
  }
}
```

## Структура проекта для Vercel

```
/
├── api/
│   └── index.go          # Основная Vercel функция
├── pkg/                  # Пакеты приложения
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── services/
├── main.go              # Локальная разработка
├── go.mod
├── go.sum
└── vercel.json          # Конфигурация Vercel
```

## Обновления для RedAPI

Файл `api/index.go` обновлен с новыми маршрутами для RedAPI:

```go
// Торренты
api.HandleFunc("/torrents/search/{imdbId}", torrentsHandler.SearchTorrents).Methods("GET")
api.HandleFunc("/torrents/movies", torrentsHandler.SearchMovies).Methods("GET")
api.HandleFunc("/torrents/series", torrentsHandler.SearchSeries).Methods("GET")
api.HandleFunc("/torrents/anime", torrentsHandler.SearchAnime).Methods("GET")
api.HandleFunc("/torrents/seasons", torrentsHandler.GetAvailableSeasons).Methods("GET")
api.HandleFunc("/torrents/search", torrentsHandler.SearchByQuery).Methods("GET")
```

## Переменные окружения для Vercel

Необходимо установить в Vercel Dashboard:

```bash
# База данных
MONGO_URI=mongodb+srv://...

# TMDB API
TMDB_ACCESS_TOKEN=your_tmdb_token

# JWT
JWT_SECRET=your_jwt_secret

# Email (если используется)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email
SMTP_PASSWORD=your_password
```

## Команды для деплоя

### 1. Через Vercel CLI:
```bash
# Установка CLI
npm i -g vercel

# Вход в аккаунт
vercel login

# Деплой
vercel
```

### 2. Через GitHub Integration:
1. Подключите репозиторий к Vercel
2. Установите переменные окружения
3. Деплой произойдет автоматически при push

## Эндпоинты API

После деплоя доступны следующие эндпоинты:

### Торренты (RedAPI):
- `GET /api/v1/torrents/search/{imdbId}` - поиск по IMDB ID
- `GET /api/v1/torrents/movies` - поиск фильмов
- `GET /api/v1/torrents/series` - поиск сериалов  
- `GET /api/v1/torrents/anime` - поиск аниме
- `GET /api/v1/torrents/seasons` - доступные сезоны
- `GET /api/v1/torrents/search` - универсальный поиск

### Фильмы и сериалы:
- `GET /api/v1/movies/*` - эндпоинты для фильмов
- `GET /api/v1/tv/*` - эндпоинты для сериалов

### Другие:
- `GET /api/v1/health` - проверка здоровья API
- `GET /docs` - документация API

## Проверка работы

После деплоя проверьте:

```bash
# Проверка здоровья
curl https://your-domain.vercel.app/api/v1/health

# Тест поиска торрентов
curl "https://your-domain.vercel.app/api/v1/torrents/search/tt0133093?type=movie"
```

## Локальная разработка

Для локальной разработки используйте:

```bash
# Запуск локального сервера
go run main.go

# Или скомпилируйте
go build -o neomovies-api .
./neomovies-api
```

## Troubleshooting

### Если возникают ошибки компиляции:
```bash
go mod tidy
go build ./api/index.go
```

### Если таймауты функций:
Увеличьте `maxDuration` в `vercel.json` (максимум 10 секунд для бесплатного плана).

### Если проблемы с CORS:
Проверьте настройки CORS в `api/index.go` - они уже настроены для разрешения всех источников.