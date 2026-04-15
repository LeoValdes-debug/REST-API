# booking-service

REST API сервис бронирования тренировок на Go.

## Стек

- **Go 1.22**
- **PostgreSQL 15** — хранение данных
- **chi** — HTTP-роутер
- **pgx** — драйвер PostgreSQL с пулом соединений
- **Docker / Docker Compose** — контейнеризация

## Архитектура

```
cmd/server/         — точка входа, сборка зависимостей, запуск сервера
internal/
  model/            — структуры данных
  repository/       — слой работы с БД (SQL-запросы)
  service/          — бизнес-логика и валидация
  handler/          — HTTP-хендлеры, роутинг
migrations/         — SQL-миграции
```

Три слоя: **Handler → Service → Repository**. Каждый слой знает только о следующем.

## Запуск

```bash
# Клонировать репо
git clone https://github.com/leovaldes-debug/booking-service
cd booking-service

# Запустить PostgreSQL + сервис
docker-compose up --build
```

Сервис поднимется на `http://localhost:8080`.

## API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | /bookings | Создать бронирование |
| GET | /bookings/{id} | Получить бронирование |
| GET | /bookings/user/{userID} | Список бронирований пользователя |
| DELETE | /bookings/{id} | Отменить бронирование |

### Примеры

**Создать бронирование**
```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "train_id": 42}'
```

```json
{
  "id": 1,
  "user_id": 1,
  "train_id": 42,
  "status": "pending",
  "created_at": "2026-04-15T10:00:00Z"
}
```

**Получить бронирование**
```bash
curl http://localhost:8080/bookings/1
```

**Список бронирований пользователя**
```bash
curl http://localhost:8080/bookings/user/1
```

**Отменить бронирование**
```bash
curl -X DELETE http://localhost:8080/bookings/1
```

## Особенности реализации

- **Graceful shutdown** - при остановке сервис дожидается завершения текущих запросов
- **Connection pooling** - pgxpool держит готовые соединения к БД
- **Индекс** по `user_id` - запросы списка бронирований работают быстро
- **Multi-stage Docker build** - итоговый образ ~10MB
