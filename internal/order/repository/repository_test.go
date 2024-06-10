package repository_test

import (
	"context"
	"testing"

	cacheMocks "github.com/deadshvt/nats-streaming-service/internal/cache/mocks"
	dbMocks "github.com/deadshvt/nats-streaming-service/internal/database/mocks"
	"github.com/deadshvt/nats-streaming-service/internal/errs"
	generator "github.com/deadshvt/nats-streaming-service/internal/generator/order"
	"github.com/deadshvt/nats-streaming-service/internal/order/repository"

	"github.com/rs/zerolog"
	a "github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// CreateOrder

// Negative tests

func TestCreateOrder_DuplicateOrder(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	order := generator.GenerateOrder()

	cache.EXPECT().CreateOrder(context.Background(), order).Return(errs.ErrOrderExists)

	err := r.CreateOrder(context.Background(), order)

	assert.ErrorIs(err, errs.ErrOrderExists)
}

// Positive tests

func TestCreateOrder_Valid(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	order := generator.GenerateOrder()

	cache.EXPECT().CreateOrder(context.Background(), order).Return(nil)
	db.EXPECT().CreateOrder(context.Background(), order).Return(nil)

	err := r.CreateOrder(context.Background(), order)

	assert.NoError(err)
}

// GetOrderByID

// Negative tests

func TestGetOrderByID_OrderNotFound(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	id := "123"

	cache.EXPECT().GetOrderByID(context.Background(), id).Return(nil, errs.ErrOrderNotFound)

	_, err := r.GetOrderByID(context.Background(), id)

	assert.ErrorIs(err, errs.ErrOrderNotFound)
}

// Positive tests

func TestGetOrderByID_Valid(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	cache := cacheMocks.NewMockOrderCache(ctrl)
	db := dbMocks.NewMockOrderDB(ctrl)
	logger := zerolog.Nop()

	r := repository.NewOrderRepository(db, cache, logger)

	id := "123"

	order := generator.GenerateOrder()
	order.OrderUid = id

	cache.EXPECT().GetOrderByID(context.Background(), id).Return(order, nil)

	o, err := r.GetOrderByID(context.Background(), id)

	assert.Equal(order, o)

	assert.NoError(err)
}
