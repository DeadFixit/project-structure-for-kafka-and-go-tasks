# Задание 3 — Go сервис с аналитикой, Redis и Kubernetes

Сервис принимает поток метрик, рассчитывает rolling average (окно 50) и z-score (>2σ = аномалия), кэширует снимки в Redis и экспортирует метрики Prometheus. Предназначен для запуска в Kubernetes (Minikube/Kind/Killercoda/Yandex.Cloud).

## Быстрый старт (PowerShell 7.5.4)
```powershell
cd task3-go-k8s-analytics
$Env:GO111MODULE="on"
# Локальный запуск (Redis должен слушать localhost:6379)
go run ./cmd/server

# Docker build
 docker build -t go-analytics:local .

# Minikube: загрузить образ
 minikube image load go-analytics:local
 kubectl apply -f deploy/k8s/redis.yaml
 kubectl apply -f deploy/k8s/configmap.yaml
 kubectl apply -f deploy/k8s/deployment.yaml
 kubectl apply -f deploy/k8s/service.yaml
 kubectl apply -f deploy/k8s/hpa.yaml
```

## HTTP API
- `POST /ingest` — `{ "timestamp": "2024-01-01T00:00:00Z", "cpu": 0.5, "rps": 1200, "latency_ms": 8 }`
- `GET /analytics` — текущий снимок (rolling avg + z-score, аномалия)
- `GET /metrics` — метрики Prometheus
- `GET /healthz` — здоровье

## Нагрузочное тестирование
Пример Locust/ab:
```powershell
# Apache Bench
ab -n 5000 -c 200 -p payload.json -T application/json http://<host>:8081/ingest
```
Цель: >1000 RPS, latency <50 мс, авто-масштабирование HPA до 4 реплик.

## Kubernetes
- **deploy/k8s/deployment.yaml** — Deployment (2 реплики, образ `go-analytics`)
- **deploy/k8s/service.yaml** — Service (ClusterIP 8081)
- **deploy/k8s/redis.yaml** — Redis (bitnami lightweight)
- **deploy/k8s/hpa.yaml** — HPA (CPU >70%, 2-5 реплик)
- **deploy/k8s/configmap.yaml** — переменные окружения