package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/deadshvt/nats-streaming-service/internal/cache"
	"github.com/deadshvt/nats-streaming-service/internal/config"
	"github.com/deadshvt/nats-streaming-service/internal/database"
	"github.com/deadshvt/nats-streaming-service/internal/middleware"
	"github.com/deadshvt/nats-streaming-service/internal/nats"
	"github.com/deadshvt/nats-streaming-service/internal/order/handler"
	"github.com/deadshvt/nats-streaming-service/internal/order/repository"
	"github.com/deadshvt/nats-streaming-service/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
)

func main() {
	baseLogger, err := logger.Init()
	if err != nil {
		log.Fatal().Msgf("Failed to initialize logger: %v", err)
	}

	config.Load(".env")

	baseLogger.Info().Msg("Loaded .env file")

	nc, err := nats.Init(os.Getenv("NATS_CLUSTER_ID"), "consumer")
	if err != nil {
		baseLogger.Fatal().Msgf("Consumer failed to connect to NATS Streaming Server: %v", err)
	}
	defer func(nc stan.Conn) {
		if err = nc.Close(); err != nil {
			baseLogger.Fatal().Msgf("Failed to close connection: %v", err)
		}
		baseLogger.Info().Msg("Consumer disconnected from NATS Streaming Server")
	}(nc)

	baseLogger.Info().Msg("Consumer connected to NATS Streaming Server")

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

	err = orderRepository.LoadCacheFromDB()
	if err != nil {
		baseLogger.Fatal().Msgf("Failed to load cache from db: %v", err)
	}

	baseLogger.Info().Msg("Consumer loaded cache from db")

	// handler
	orderHandlerLogger := logger.NewLogger(baseLogger, "OrderHandler")
	orderHandler := handler.NewOrderHandler(orderRepository, orderHandlerLogger)

	if _, err = nc.Subscribe(os.Getenv("NATS_SUBJECT"), func(m *stan.Msg) {
		err = orderHandler.CreateOrder(m.Data)
		if err != nil {
			baseLogger.Error().Msgf("Failed to handle order: %v", err)
		}
	}); err != nil {
		baseLogger.Fatal().Msgf("Consumer failed to subscribe: %v", err)
	}

	baseLogger.Info().Msg("Consumer subscribed to NATS Streaming Server")

	r := mux.NewRouter()

	r.HandleFunc("/", orderHandler.GetOrderID).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/order/{id}", orderHandler.GetOrderByID).Methods(http.MethodGet)

	mx := middleware.Panic(baseLogger, r)
	mx = middleware.Logging(baseLogger, mx)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	baseLogger.Info().Msgf("Starting server on port=:%s", port)

	if err = http.ListenAndServe(":"+port, mx); err != nil {
		baseLogger.Fatal().Msgf("Failed to start server: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	baseLogger.Info().Msg("Shutting down consumer...")
}
