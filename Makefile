.PHONY: prometheus migrate-up migrate-down lint consumer producer server container loadtest

PROMETHEUS_CONFIG=prometheus.yml

DB=postgresql
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=order
DB_SSLMODE=disable

prometheus:
	prometheus --config.file=$(PROMETHEUS_CONFIG)

migrate-up:
	migrate -path internal/database/migration -database "$(DB)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose up

migrate-down:
	migrate -path internal/database/migration -database "$(DB)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose down

lint:
	golangci-lint run

consumer:
	go run cmd/consumer/main.go

producer:
	go run cmd/producer/main.go

server:
	go run cmd/server/main.go

container:
	docker compose up -d

loadtest:
	go run loadtest/loadtest.go
