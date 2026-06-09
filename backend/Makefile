.PHONY: run test db-up db-down migrate

run:
	go run ./cmd/api

test:
	go test ./...

db-up:
	docker compose -f deployments/local/docker-compose.yml up -d postgres redis

db-down:
	docker compose -f deployments/local/docker-compose.yml down

migrate:
	docker compose -f deployments/local/docker-compose.yml run --rm flyway
