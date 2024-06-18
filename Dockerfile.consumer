FROM golang:alpine as builder

WORKDIR /project

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/consumer ./cmd/consumer

FROM alpine:latest

WORKDIR /project

COPY .env .env
COPY ./web ./web
COPY --from=builder /project/bin/consumer /project/bin/consumer

CMD ["/project/bin/consumer"]
