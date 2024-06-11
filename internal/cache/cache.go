package cache

import (
	"context"

	"github.com/deadshvt/nats-streaming-service/internal/cache/in_memory"
	"github.com/deadshvt/nats-streaming-service/internal/entity"
	"github.com/deadshvt/nats-streaming-service/internal/errs"
)

const (
	InMemory = "in_memory"
)

type OrderCache interface {
	GetOrderByID(ctx context.Context, id string) (*entity.Order, error)
	CreateOrder(ctx context.Context, order *entity.Order) error
}

func InitOrderCache(cacheType string) (OrderCache, error) {
	switch cacheType {
	case InMemory:
		return in_memory.NewInMemory(), nil
	default:
		return nil, errs.ErrUnsupportedCacheType
	}
}
