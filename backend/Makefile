COMPOSE := docker compose --env-file .env -f deployments/local/docker-compose.yml

.PHONY: run test db-up db-down migrate

run:
	go run ./cmd/api

test:
	go test ./...

db-up:
	$(COMPOSE) up -d postgres redis

db-down:
	$(COMPOSE) down

migrate:
	$(COMPOSE) run --rm flyway
