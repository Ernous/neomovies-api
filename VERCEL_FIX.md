# Исправление ошибки деплоя на Vercel

## ❌ Ошибка
```
index.go:12:2: use of internal package neomovies-api/internal/config not allowed
```

## 🔍 Причина
Go не разрешает импорт пакетов из папки `internal` для кода, находящегося вне корневого модуля. Файл `api/index.go` не может импортировать `neomovies-api/internal/*` пакеты.

## ✅ Решение
Переименована папка `internal` в `pkg` и обновлены все импорты:

### Изменения в структуре:
```
internal/  →  pkg/
├── config/
├── database/  
├── handlers/
├── middleware/
├── models/
└── services/
```

### Обновленные импорты:
```go
// Было:
"neomovies-api/internal/config"
"neomovies-api/internal/database"
"neomovies-api/internal/handlers"
"neomovies-api/internal/middleware"
"neomovies-api/internal/services"

// Стало:
"neomovies-api/pkg/config"
"neomovies-api/pkg/database"
"neomovies-api/pkg/handlers"
"neomovies-api/pkg/middleware"
"neomovies-api/pkg/services"
```

## 📁 Затронутые файлы:
- `api/index.go` - основная Vercel функция
- `main.go` - локальная разработка
- Все файлы в `pkg/` - обновлены импорты

## 🧪 Проверка
- ✅ `go build ./api/index.go` - компилируется
- ✅ `go build .` - основной проект компилируется
- ✅ Все функции RedAPI работают корректно
- ✅ Тесты пройдены

## 🚀 Готово к деплою
Теперь проект должен успешно деплоиться на Vercel без ошибок импорта internal пакетов.

### Команда для деплоя:
```bash
vercel --prod
```

### Проверка после деплоя:
```bash
curl https://your-domain.vercel.app/api/v1/health
curl "https://your-domain.vercel.app/api/v1/torrents/search/tt0133093?type=movie"
```