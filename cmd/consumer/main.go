package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deadshvt/nats-streaming-service/internal/cache"
	"github.com/deadshvt/nats-streaming-service/internal/config"
	"github.com/deadshvt/nats-streaming-service/internal/database"
	"github.com/deadshvt/nats-streaming-service/internal/middleware"
	"github.com/deadshvt/nats-streaming-service/internal/nats"
	"github.com/deadshvt/nats-streaming-service/internal/order/handler"
	"github.com/deadshvt/nats-streaming-service/internal/order/repository"
	"github.com/deadshvt/nats-streaming-service/internal/prometheus"
	"github.com/deadshvt/nats-streaming-service/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	baseLogger, err := logger.Init()
	if err != nil {
		log.Fatal().Msgf("Failed to initialize logger: %v", err)
	}

	baseLogger.Info().Msg("Initialized logger")

	config.Load(".env")

	baseLogger.Info().Msg("Loaded .env file")

	sc, err := nats.Init(os.Getenv("NATS_CLUSTER_ID"), "consumer", os.Getenv("NATS_URL"))
	if err != nil {
		baseLogger.Fatal().Msgf("Consumer failed to connect to NATS Streaming Server: %v", err)
	}
	defer func(nc stan.Conn) {
		if err = nc.Close(); err != nil {
			baseLogger.Fatal().Msgf("Failed to close connection: %v", err)
		}
		baseLogger.Info().Msg("Consumer disconnected from NATS Streaming Server")
	}(sc)

	baseLogger.Info().Msg("Consumer connected to NATS Streaming Server")

	ctx := context.Background()

	// db
	orderDB, err := database.InitOrderDB(os.Getenv("DB_TYPE"))
	if err != nil {
		baseLogger.Fatal().Msgf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err = orderDB.Disconnect(); err != nil {
			baseLogger.Fatal().Msgf("Failed to close database: %v", err)
		}
		baseLogger.Info().Msg("Consumer disconnected from database")
	}()

	baseLogger.Info().Msg("Consumer connected to database")

	// cache
	orderCache, err := cache.InitOrderCache(os.Getenv("CACHE_TYPE"))

	// repository
	orderRepositoryLogger := logger.NewLogger(baseLogger, "OrderRepository")
	orderRepository := repository.NewOrderRepository(orderDB, orderCache, orderRepositoryLogger)

	err = orderRepository.LoadCacheFromDB(ctx)
	if err != nil {
		baseLogger.Fatal().Msgf("Failed to load cache from db: %v", err)
	}

	baseLogger.Info().Msg("Consumer loaded cache from db")

	// handler
	orderHandlerLogger := logger.NewLogger(baseLogger, "OrderHandler")
	orderHandler := handler.NewOrderHandler(orderRepository, orderHandlerLogger)

	if _, err = sc.Subscribe(os.Getenv("NATS_SUBJECT"), func(m *stan.Msg) {
		err = orderHandler.CreateOrder(ctx, m.Data)
		if err != nil {
			baseLogger.Error().Msgf("Failed to handle order: %v", err)
		}
	}); err != nil {
		baseLogger.Fatal().Msgf("Consumer failed to subscribe: %v", err)
	}

	baseLogger.Info().Msg("Consumer subscribed to NATS Streaming Server")

	metrics := prometheus.NewMetrics()
	metrics.Register()

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	r.HandleFunc("/", metrics.Handler("/", orderHandler.GetOrderID)).
		Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/order/{id}", metrics.Handler("/order/{id}", orderHandler.GetOrderByID)).
		Methods(http.MethodGet)

	mx := middleware.Panic(baseLogger, r)
	mx = middleware.Logging(baseLogger, mx)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	baseLogger.Info().Msgf("Starting server on port=:%s", port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mx,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			baseLogger.Fatal().Msgf("Failed to start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	baseLogger.Info().Msg("Shutting down consumer...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		baseLogger.Fatal().Msgf("Failed to gracefully shutdown server: %v", err)
	}

	baseLogger.Info().Msg("Server gracefully stopped")
}
