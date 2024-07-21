package repository_test

import (
	"context"
	"testing"

	"github.com/deadshvt/nats-streaming-service/internal/errs"
	orderGenerator "github.com/deadshvt/nats-streaming-service/internal/generator/order"
	"github.com/deadshvt/nats-streaming-service/internal/order/repository"
	cacheMocks "github.com/deadshvt/nats-streaming-service/internal/order/repository/cache/mocks"
	dbMocks "github.com/deadshvt/nats-streaming-service/internal/order/repository/database/mocks"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// CreateOrder

// Negative tests

func TestCreateOrder_DuplicateOrder(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	order := orderGenerator.RandomOrder()

	cache.EXPECT().CreateOrder(context.Background(), order).Return(errs.ErrOrderExists)

	err := r.CreateOrder(context.Background(), order)

	a.ErrorIs(err, errs.ErrOrderExists)
}

// Positive tests

func TestCreateOrder_Valid(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	order := orderGenerator.RandomOrder()

	cache.EXPECT().CreateOrder(context.Background(), order).Return(nil)
	db.EXPECT().CreateOrder(context.Background(), order).Return(nil)

	err := r.CreateOrder(context.Background(), order)

	a.NoError(err)
}

// GetOrderByID

// Negative tests

func TestGetOrderByID_OrderNotFound(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	id := "123"

	cache.EXPECT().GetOrderByID(context.Background(), id).Return(nil, errs.ErrOrderNotFound)

	_, err := r.GetOrderByID(context.Background(), id)

	a.ErrorIs(err, errs.ErrOrderNotFound)
}

// Positive tests

func TestGetOrderByID_Valid(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	id := "123"

	order := orderGenerator.RandomOrder()
	order.OrderUid = id

	cache.EXPECT().GetOrderByID(context.Background(), id).Return(order, nil)

	o, err := r.GetOrderByID(context.Background(), id)

	a.Equal(order, o)

	a.NoError(err)
}
