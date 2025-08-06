# Neo Movies API (Go)

Современный API для поиска фильмов и сериалов с интеграцией TMDB, написанный на Go.

## 🚀 Особенности

- **Go Backend**: Высокопроизводительный backend на Go
- **TMDB Integration**: Полная интеграция с The Movie Database API
- **JWT Authentication**: Безопасная аутентификация пользователей
- **MongoDB**: Надежное хранение пользовательских данных
- **Vercel Ready**: Готов к деплою на Vercel
- **Beautiful Docs**: Современная документация API с Scalar
- **CORS Support**: Поддержка кросс-доменных запросов

## 📋 API Documentation

После запуска сервера документация доступна по адресу:
- **Local**: http://localhost:3000/docs
- **Production**: https://your-app.vercel.app/docs

## 🛠 Технологии

- **Go 1.21+**
- **Gorilla Mux** - HTTP маршрутизация
- **MongoDB** - База данных
- **JWT** - Аутентификация
- **TMDB API** - Данные о фильмах
- **Scalar** - Документация API

## 🏃‍♂️ Быстрый старт

### Локальная разработка

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/your-username/neomovies-api-go.git
cd neomovies-api-go
```

2. **Настройте переменные окружения:**
```bash
cp .env.example .env
# Отредактируйте .env файл с вашими настройками
```

3. **Установите зависимости:**
```bash
go mod download
```

4. **Запустите сервер:**
```bash
go run main.go
```

Сервер будет доступен по адресу: http://localhost:3000

### Деплой на Vercel

1. **Подключите репозиторий к Vercel**

2. **Настройте переменные окружения в Vercel:**
   - `MONGO_URI` - Строка подключения к MongoDB Atlas
   - `TMDB_API_KEY` - Ключ API от TMDB
   - `JWT_SECRET` - Секретный ключ для JWT токенов

3. **Деплой произойдет автоматически**

## 🔑 Получение TMDB API Key

1. Зарегистрируйтесь на [TMDB](https://www.themoviedb.org/)
2. Перейдите в [Settings > API](https://www.themoviedb.org/settings/api)
3. Создайте новый API ключ
4. Добавьте ключ в переменную окружения `TMDB_API_KEY`

## 📚 API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Авторизация пользователя
- `GET /api/v1/auth/profile` - Получить профиль пользователя (требует JWT)
- `PUT /api/v1/auth/profile` - Обновить профиль пользователя (требует JWT)

### Фильмы
- `GET /api/v1/movies/search` - Поиск фильмов
- `GET /api/v1/movies/popular` - Популярные фильмы
- `GET /api/v1/movies/top-rated` - Высокорейтинговые фильмы
- `GET /api/v1/movies/upcoming` - Предстоящие фильмы
- `GET /api/v1/movies/now-playing` - Фильмы в прокате
- `GET /api/v1/movies/{id}` - Информация о фильме
- `GET /api/v1/movies/{id}/recommendations` - Рекомендации
- `GET /api/v1/movies/{id}/similar` - Похожие фильмы

### Сериалы
- `GET /api/v1/tv/search` - Поиск сериалов
- `GET /api/v1/tv/popular` - Популярные сериалы
- `GET /api/v1/tv/top-rated` - Высокорейтинговые сериалы
- `GET /api/v1/tv/on-the-air` - Сериалы в эфире
- `GET /api/v1/tv/airing-today` - Сериалы сегодня
- `GET /api/v1/tv/{id}` - Информация о сериале
- `GET /api/v1/tv/{id}/recommendations` - Рекомендации
- `GET /api/v1/tv/{id}/similar` - Похожие сериалы

### Избранное (требует авторизации)
- `GET /api/v1/favorites` - Список избранных фильмов
- `POST /api/v1/favorites/{id}` - Добавить в избранное
- `DELETE /api/v1/favorites/{id}` - Удалить из избранного

### Служебные
- `GET /api/v1/health` - Проверка работоспособности API

## 🔐 Аутентификация

API использует JWT токены для аутентификации. Для доступа к защищенным эндпоинтам добавьте заголовок:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

## 📝 Примеры использования

### Регистрация пользователя
```bash
curl -X POST https://your-app.vercel.app/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "Иван Иванов"
  }'
```

### Поиск фильмов
```bash
curl "https://your-app.vercel.app/api/v1/movies/search?query=Мстители&language=ru-RU"
```

### Добавление в избранное
```bash
curl -X POST https://your-app.vercel.app/api/v1/favorites/550 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🌟 Новые возможности по сравнению с JS версией

- **Улучшенная производительность** благодаря Go
- **Современная документация** с Scalar вместо Swagger UI
- **Лучшая обработка ошибок** и валидация данных
- **Типобезопасность** и строгая типизация
- **Оптимизированное подключение к БД** с пулом соединений
- **Улучшенная безопасность** JWT токенов

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для фичи (`git checkout -b feature/AmazingFeature`)
3. Зафиксируйте изменения (`git commit -m 'Add some AmazingFeature'`)
4. Отправьте в ветку (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🙏 Благодарности

- [TMDB](https://www.themoviedb.org/) за предоставление API данных о фильмах
- [Scalar](https://scalar.com/) за красивую документацию API
- [Vercel](https://vercel.com/) за хостинг

---

**Сделано с ❤️ на Go**
