package in_memory

import (
	"context"
	"sync"

	"github.com/deadshvt/nats-streaming-service/internal/entity"
	"github.com/deadshvt/nats-streaming-service/internal/errs"
)

type InMemory struct {
	orders map[string]*entity.Order
	mu     *sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		orders: make(map[string]*entity.Order),
		mu:     &sync.RWMutex{},
	}
}

func (r *InMemory) CreateOrder(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, ok := r.orders[order.OrderUid]; ok {
			return errs.ErrOrderExists
		}
	}

	r.orders[order.OrderUid] = order

	return nil
}

func (r *InMemory) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if _, ok := r.orders[id]; !ok {
			return nil, errs.ErrOrderNotFound
		}
	}

	return r.orders[id], nil
}
