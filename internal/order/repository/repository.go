package repository

import (
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

func NewOrderRepository(db database.OrderDB, cache cache.OrderCache, logger zerolog.Logger) *OrderRepository {
	return &OrderRepository{
		DB:     db,
		Cache:  cache,
		Logger: logger,
	}
}

func (r *OrderRepository) CreateOrder(order *entity.Order) error {
	r.Logger.Info().Msg("Creating order...")

	err := r.Cache.CreateOrder(order)
	if err != nil {
		return err
	}

	err = r.DB.CreateOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetOrderByID(id string) (*entity.Order, error) {
	r.Logger.Info().Msg("Getting order by id...")

	order, err := r.Cache.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) LoadCacheFromDB() error {
	r.Logger.Info().Msg("Loading cache from db...")

	orders, err := r.DB.GetAllOrders()
	if err != nil {
		return err
	}

	for i := range orders {
		err = r.Cache.CreateOrder(orders[i])
		if err != nil {
			return err
		}
	}

	return nil
}
