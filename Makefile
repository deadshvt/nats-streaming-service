.PHONY: run wait-postgres migrate-up migrate-down lint consumer producer containers-up containers-down loadtest test

DB=postgresql
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=order
DB_SSLMODE=disable
DB_DSN="$(DB)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

run: containers-up migrate-up

wait-postgres:
	@until docker exec -t $(shell docker compose ps -q postgres) pg_isready -U postgres; do sleep 1; done

migrate-up: wait-postgres
	docker run --rm --network $(shell docker inspect --format='{{.HostConfig.NetworkMode}}' $(shell docker compose ps -q postgres)) -v $(shell pwd)/internal/database/migration:/migrations migrate/migrate -path /migrations/ -database $(DB_DSN) -verbose up

migrate-down: wait-postgres
	docker run --rm --network $(shell docker inspect --format='{{.HostConfig.NetworkMode}}' $(shell docker compose ps -q postgres)) -v $(shell pwd)/internal/database/migration:/migrations migrate/migrate -path /migrations/ -database $(DB_DSN) -verbose down

lint:
	golangci-lint run

consumer:
	docker compose up consumer

producer:
	docker compose up producer

containers-up:
	docker compose up -d postgres nats-streaming prometheus

containers-down:
	docker compose down

loadtest:
	go run loadtest/loadtest.go

test:
	go test -v ./...
