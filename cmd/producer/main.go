package main

import (
	"encoding/json"
	"github.com/deadshvt/nats-streaming-service/internal/nats"
	"os"
	"strings"
	"time"

	"github.com/deadshvt/nats-streaming-service/internal/config"
	"github.com/deadshvt/nats-streaming-service/internal/entity"
	generator "github.com/deadshvt/nats-streaming-service/internal/generator/order"
	"github.com/deadshvt/nats-streaming-service/pkg/logger"

	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
)

const (
	Count = 5
	Dir   = "./schema"
)

func main() {
	baseLogger, err := logger.Init()
	if err != nil {
		log.Fatal().Msgf("Failed to initialize logger: %v", err)
	}

	config.Load(".env")

	sc, err := nats.Init(os.Getenv("NATS_CLUSTER_ID"), "producer", os.Getenv("NATS_URL"))
	if err != nil {
		panic(err)
	}
	defer func(sc stan.Conn) {
		if err = sc.Close(); err != nil {
			baseLogger.Error().Msgf("Failed to close connection: %v", err)
		}
		baseLogger.Info().Msg("Producer disconnected from NATS Streaming Server")
	}(sc)

	baseLogger.Info().Msg("Producer connected to NATS Streaming Server")

	files, err := os.ReadDir(Dir)
	if err != nil {
		baseLogger.Error().Msgf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		fileData, err := os.ReadFile(Dir + "/" + file.Name())
		if err != nil {
			baseLogger.Error().Msgf("Failed to read file: %v", err)
			continue
		}
		if err = sc.Publish(os.Getenv("NATS_SUBJECT"), fileData); err != nil {
			baseLogger.Error().Msgf("Failed to publish message: %v", err)
			continue
		}

		var order entity.Order
		err = json.Unmarshal(fileData, &order)
		if err != nil {
			baseLogger.Error().Msgf("Failed to unmarshal order: %v", err)
		}

		logger.LogWithParams(baseLogger, "Published order", struct {
			OrderID string
		}{OrderID: order.OrderUid})

		time.Sleep(1 * time.Second)
	}

	for i := 0; i < Count; i++ {
		order := generator.GenerateOrder()
		data, err := json.Marshal(*order)
		if err != nil {
			baseLogger.Error().Msgf("Failed to marshal order: %v", err)
			continue
		}
		if err = sc.Publish(os.Getenv("NATS_SUBJECT"), data); err != nil {
			baseLogger.Error().Msgf("Failed to publish message: %v", err)
			continue
		}

		logger.LogWithParams(baseLogger, "Published order", struct {
			OrderID string
		}{OrderID: order.OrderUid})

		time.Sleep(2 * time.Second)
	}
}
