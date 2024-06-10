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

1. **Set up migrations**

```shell
make migrate-up
```

2. **Set up containers**

```shell
make containers
```

3. **Run the consumer:**

```shell
make consumer
```

4. **Run the producer:**

```shell
make producer
```

5. **Open browser in `http://localhost`**

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

## Monitoring

1. **Set up Prometheus**

```shell
make prometheus
```

2. **Open browser in `http://localhost:9090`**
