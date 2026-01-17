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
"username": "trump",
"password": "111222333",
"email": "trump@example.com"
}'

curl -X GET http://localhost:8080/api/v1/login/testuser
curl -X GET http://localhost:8080/api/v1/login/trump/111222333

curl -X POST http://localhost:8080/api/v1/pst/write \
-H "Content-Type: application/json" \
-H "token: kdiejfji393874" \
-d '{
"title": "Alice巡游记",
"content": "Alice巡游记，这是第一篇博客",
"userid": 1
}'


curl -X POST http://localhost:8080/api/v1/pst/write \
-H "Content-Type: application/json" \
-H "token: kdiejfji393874" \
-d '{
"title": "Alice巡游记第二章",
"content": "Alice巡游记，这是第二篇博客",
"userid": 1
}'

curl -X POST http://localhost:8080/api/v1/pst/write \
-H "Content-Type: application/json" \
-d '{
"userid": 1,
"page": "111222333",
"size": "trump@example.com"
}'

curl -X GET http://localhost:8080/api/v1/pst/1/2/5
curl -X GET http://localhost:8080/api/v1/pst/1/2/5 \
-H "token: kdiejfji393874"