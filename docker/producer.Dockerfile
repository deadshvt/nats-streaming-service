FROM golang:alpine as builder

WORKDIR /project

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/producer ./cmd/producer

FROM alpine:latest

WORKDIR /project

COPY .env .env
COPY ./schema ./schema
COPY --from=builder /project/bin/producer /project/bin/producer

CMD ["/project/bin/producer"]
