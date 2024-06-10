package repository

import (
	"context"

	"github.com/deadshvt/nats-streaming-service/internal/cache"
	"github.com/deadshvt/nats-streaming-service/internal/database"
	"github.com/deadshvt/nats-streaming-service/internal/entity"

	"github.com/rs/zerolog"
)

type OrderRepository struct {
	DB     database.OrderDB
	Cache  cache.OrderCache
	Logger zerolog.Logger
}

func NewOrderRepository(db database.OrderDB, cache cache.OrderCache, logger zerolog.Logger) entity.OrderRepository {
	return &OrderRepository{
		DB:     db,
		Cache:  cache,
		Logger: logger,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	r.Logger.Info().Msg("Creating order...")

	err := r.Cache.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	err = r.DB.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	r.Logger.Info().Msg("Getting order by id...")

	order, err := r.Cache.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) LoadCacheFromDB(ctx context.Context) error {
	r.Logger.Info().Msg("Loading cache from db...")

	orders, err := r.DB.GetAllOrders(ctx)
	if err != nil {
		return err
	}

	for i := range orders {
		err = r.Cache.CreateOrder(ctx, orders[i])
		if err != nil {
			return err
		}
	}

	return nil
}
