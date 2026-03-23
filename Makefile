.PHONY: build up down test build-run-app migrate-up migrate-down migrate-down-all seed-all db-shell logs app-logs

build:
	docker compose -f docker-compose.yml build

up:
	docker compose -f docker-compose.yml up -d

build-run-app:
	docker compose -f docker-compose.yml up --build -d app

down:
	docker compose -f docker-compose.yml down

test:
	go test ./...

migrate-up:
	docker compose -f docker-compose.yml run --rm migrations

migrate-down:
	docker compose -f docker-compose.yml run --rm migrations /migrate down

migrate-down-all:
	docker compose -f docker-compose.yml run --rm migrations /migrate --down --down-all

seed-all:
	docker compose -f docker-compose.yml run --rm migrations /migrate --seed

db-shell:
	docker exec -it sports-calendar-db bash

logs:
	docker compose -f docker-compose.yml logs -f

app-logs:
	docker compose -f docker-compose.yml logs -f app