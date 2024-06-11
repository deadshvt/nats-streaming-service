# nats-streaming-service

## Overview
This project is a NATS Streaming Service that processes orders. 
It includes a consumer to process messages, a producer to publish messages, and a server to expose HTTP endpoints for interacting with orders. 
Prometheus is used for monitoring, and migration scripts manage the database schema.

## Installation

1. **Clone the repository:**

```shell
git clone https://github.com/deadshvt/nats-streaming-service.git
```

2. **Go to the project directory:**

```shell
cd nats-streaming-service
```

3. **Install dependencies:**

```shell
go mod tidy
```

## Running the application

1. **Set up nats, postgres, prometheus and migrations**

```shell
make run
```

2. **Run the consumer:**

```shell
make consumer
```

3. **Run the producer:**

```shell
make producer
```

4. **Open browser in `http://localhost:8080` to see home page to get order**

5. **Open browser in `http://localhost:9090` to see the metrics**

## Running tests

1. **Unit tests:**

```shell
make test
```

2. **Loading testing:**

```shell
make loadtest
```

## Linting

```shell
make lint
```
