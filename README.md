# Subscription Service

HTTP API сервис на Go для управления пользовательскими подписками с БД PostgreSQL.

## Быстрый старт

### Требования
- Docker и Docker Compose
- Make

### Запуск

1. **Клонирование проекта**:
   ```bash
   git clone https://github.com/reedtzdrvn/subscription-service
   ```

2. **Запуск сервисов**:
   ```bash
   make
   ```

### Альтернативный способ: Запуск без Docker

1. **Запуск сервера напрямую**:
   ```bash
   go run cmd/server/main.go
   ```
   
Убедитесь, что PostgreSQL запущен и настроена конфигурация окружения.

## HTTP API Маршруты

### Управление Подписками

| Метод | Конечная точка | Описание | Параметры |
|--------|----------|-------------|------------|
| `GET` | `/subscriptions` | Список подписок с опциональными фильтрами | Query: `user_id`, `service_name`, `from`, `to` |
| `POST` | `/subscriptions` | Создать новую подписку | JSON тело (см. ниже) |
| `GET` | `/subscriptions/:id` | Получить подписку по ID | Path: `id` (UUID) |
| `PUT` | `/subscriptions/:id` | Обновить подписку | Path: `id` (UUID), JSON тело |
| `DELETE` | `/subscriptions/:id` | Удалить подписку | Path: `id` (UUID) |
| `GET` | `/subscriptions/sum` | Рассчитать общую стоимость подписок | Query: `user_id`, `service_name`, `from*`, `to*` |

*Обязательные параметры

### Примеры Запросов/Ответов

#### Создание Подписки
```bash
POST /subscriptions
Content-Type: application/json

{
  "service_name": "Netflix",
  "price": 1500,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2024",
  "end_date": "12-2024"
}
```

#### Список с фильтрами
```bash
GET /subscriptions?user_id=550e8400-e29b-41d4-a716-446655440000&from=01-2024&to=12-2024
```

#### Расчет суммы
```bash
GET /subscriptions/sum?user_id=550e8400-e29b-41d4-a716-446655440000&from=01-2024&to=12-2024
```

### Формат дат
Все даты используют формат `MM-YYYY` (например, `01-2024` для января 2024).

### Формат ответа
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Netflix",
  "price": 1500,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2024",
  "end_date": "12-2024",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## Структура проекта

```
.
├── cmd/server/           # Главная точка входа приложения
├── internal/
│   ├── config/          # Управление конфигурацией
│   ├── db/              # Подключение к базе данных
│   ├── domain/          # Бизнес-сущности
│   ├── logger/          # Утилиты логирования
│   ├── repository/      # Слой доступа к данным
│   ├── transport/http/  # HTTP обработчики
│   └── usecase/         # Бизнес-логика
├── migrations/          # Миграции базы данных
├── api/                 # OpenAPI спецификация
├── docker-compose.yml   # Docker сервисы
├── Dockerfile          # Контейнер приложения
└── makefile            # Команды сборки
```

## Конфигурация

Приложение использует переменные окружения для конфигурации. Проверьте файл `.env` для доступных опций.

## База Данных

Сервис использует PostgreSQL с автоматическими миграциями через Goose. Миграции находятся в директории `migrations/`.
