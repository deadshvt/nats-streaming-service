package cache

import (
	"github.com/deadshvt/nats-streaming-service/internal/cache/in_memory"
	"github.com/deadshvt/nats-streaming-service/internal/entity"
	"github.com/deadshvt/nats-streaming-service/internal/errs"
)

type OrderCache interface {
	GetOrderByID(id string) (*entity.Order, error)
	CreateOrder(order *entity.Order) error
}

func InitOrderCache(cacheType string) (OrderCache, error) {
	switch cacheType {
	case "in_memory":
		return in_memory.NewInMemory(), nil
	default:
		return nil, errs.ErrUnsupportedCacheType
	}
}
