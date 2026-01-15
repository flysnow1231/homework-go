.PHONY: run tidy up down

tidy:
	go mod tidy

run:
	go run ./cmd/api

up:
	docker compose -f deploy/docker-compose.yaml up -d

down:
	docker compose -f deploy/docker-compose.yaml down
