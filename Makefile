.PHONY: lint consumer producer server container loadtest

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
