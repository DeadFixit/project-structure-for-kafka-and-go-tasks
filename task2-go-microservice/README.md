# Задание 2 — Go high-load микросервис

Этот проект реализует REST API для CRUD-операций над пользователями с асинхронным аудитом, метриками Prometheus, rate limiting и интеграцией с MinIO.

## Быстрый старт (PowerShell 7.5.4)
```powershell
# 1) Сборка
cd task2-go-microservice
$Env:GO111MODULE="on"
go build ./...

# 2) Запуск локально
$Env:MINIO_ENDPOINT="localhost:9000"
$Env:MINIO_ACCESS_KEY="minioadmin"
$Env:MINIO_SECRET_KEY="minioadmin"
$Env:MINIO_BUCKET="demo-bucket"

# Запуск MinIO и сервиса через docker-compose
 docker compose up --build
# или только бинарник
 go run ./...
```

## HTTP API
- `GET /api/users` — список пользователей
- `GET /api/users/{id}` — получить пользователя
- `POST /api/users` — создать (body: `{ "name": "", "email": "" }`)
- `PUT /api/users/{id}` — обновить
- `DELETE /api/users/{id}` — удалить
- `POST /api/integration/minio/upload` — записать sample-объект в MinIO
- `GET /metrics` — метрики Prometheus
- `GET /healthz` — проверка готовности

## Архитектура
- **handlers/** — HTTP-слой (gorilla/mux)
- **services/** — доменная логика (память + MinIO-клиент)
- **utils/rate_limiter.go** — rate limiter на 1000 rps (burst 5000)
- **metrics/** — Prometheus middleware (RPS + latency)

Асинхронные действия выполняются goroutine-ами в `utils/logger.go` (audit лог), а также при вызове интеграций (`handlers/integration_handler.go`).

## Нагрузочное тестирование
Пример команды wrk:
```powershell
wrk -t12 -c500 -d60s http://localhost:8080/api/users
```
Ожидания: RPS > 1000, latency < 10 мс, ошибок — 0. Фактические результаты занесите в `docs/report-template.md` и при необходимости экспортируйте в .docx/.pdf.

## Контейнеризация
- `Dockerfile` — минимальный образ на базе golang:1.22-alpine.
- `docker-compose.yml` — сервис, MinIO и Prometheus (по умолчанию слушает `:8080`).
