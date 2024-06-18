include .env
export

DB_DSN="postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

run: containers-up migrate-up ### Run containers: postgres, nats-streaming, prometheus and migrate up
.PHONY: run

stop: migrate-down ### Stop all containers and migrate down
	docker compose down
.PHONY: stop

wait-postgres: ### Wait until Postgres is ready
	@until docker exec -t $(shell docker compose ps -q postgres) pg_isready -U postgres; do sleep 1; done
.PHONY: wait-postgres

migrate-up: wait-postgres ### Migrate up
	docker run --rm --network $(shell docker inspect --format='{{.HostConfig.NetworkMode}}' $(shell docker compose ps -q postgres)) -v $(shell pwd)/migrations:/migrations migrate/migrate -path /migrations/ -database $(DB_DSN) -verbose up
.PHONY: migrate-up

migrate-down: wait-postgres ### Migrate down
	echo "y" | docker run --rm --network $(shell docker inspect --format='{{.HostConfig.NetworkMode}}' $(shell docker compose ps -q postgres)) -v $(shell pwd)/migrations:/migrations migrate/migrate -path /migrations/ -database $(DB_DSN) -verbose down -all
.PHONY: migrate-down

lint: ### Run linting
	golangci-lint run
.PHONY: lint

consumer: ### Run consumer
	docker compose up consumer
.PHONY: consumer

producer: ### Run producer
	docker compose up producer
.PHONY: producer

containers-up: ### Run containers: postgres, nats-streaming, prometheus
	docker compose up -d postgres nats-streaming prometheus
.PHONY: containers-up

containers-down: ### Stop containers: postgres, nats-streaming, prometheus
	docker compose down postgres nats-streaming prometheus
.PHONY: containers-down

loadtest: ### Run load testing
	go run loadtest/loadtest.go
.PHONY: loadtest

test: ### Run tests
	go test -v ./...
.PHONY: test
