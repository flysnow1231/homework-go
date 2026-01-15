# bangkok (Go Web Skeleton)

Gin + GORM(MySQL) + Redis + RabbitMQ + Zap + Viper

## Endpoints
- Health: GET /healthz
- Readiness: GET /readyz
- Demo API:
  - POST /api/v1/register
  - GET  /api/v1/user/:id

## Quick start (local)
```bash
cd blog
docker compose -f deploy/docker-compose.yaml up -d
go mod tidy
go run ./cmd/api
```

RabbitMQ management UI:
- http://localhost:15672 (guest/guest)

## Config
Edit `configs/config.yaml` or override with env:
- BLOG_APP_HTTP_ADDR=:8080
- BLOG_MYSQL_DSN=...

## 注册
curl -X POST http://localhost:8080/api/v1/register \
-H "Content-Type: application/json" \
-d '{
"username": "testuser",
"password": "123456",
"email": "testuser@example.com"
}'

curl -X GET http://localhost:8080/api/v1/users/testuser
