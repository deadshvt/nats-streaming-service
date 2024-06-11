package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deadshvt/nats-streaming-service/internal/entity"
	"github.com/deadshvt/nats-streaming-service/internal/errs"
	generator "github.com/deadshvt/nats-streaming-service/internal/generator/order"
	"github.com/deadshvt/nats-streaming-service/internal/order/handler"
	"github.com/deadshvt/nats-streaming-service/internal/order/mocks"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	a "github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// CreateOrder

// Negative tests

func TestCreateOrder_InvalidJSON(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockOrderRepository(ctrl)
	logger := zerolog.Nop()
	h := handler.NewOrderHandler(repo, logger)

	jsonOrder := []byte(`{[[[[[[}`)

	err := h.CreateOrder(context.Background(), jsonOrder)

	assert.ErrorContains(err, errs.ErrJSONUnmarshal.Error())
}

func TestCreateOrder_InvalidOrder(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockOrderRepository(ctrl)
	logger := zerolog.Nop()
	h := handler.NewOrderHandler(repo, logger)

	order := generator.GenerateOrder()
	order.OrderUid = ""
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}

	err = h.CreateOrder(context.Background(), jsonOrder)

	assert.ErrorContains(err, errs.ErrValidateOrder.Error())
}

func TestCreateOrder_DuplicateOrder(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockOrderRepository(ctrl)
	logger := zerolog.Nop()
	h := handler.NewOrderHandler(repo, logger)

	order := generator.GenerateOrder()
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}

	repo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(&entity.Order{})).Return(errs.ErrOrderExists)

	err = h.CreateOrder(context.Background(), jsonOrder)

	assert.ErrorContains(err, errs.ErrOrderExists.Error())
	assert.ErrorContains(err, errs.ErrCreateOrder.Error())
}

// Positive tests

func TestCreateOrder_Valid(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockOrderRepository(ctrl)
	logger := zerolog.Nop()
	h := handler.NewOrderHandler(repo, logger)

	order := generator.GenerateOrder()
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}

	repo.EXPECT().CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(&entity.Order{})).Return(nil)

	err = h.CreateOrder(context.Background(), jsonOrder)

	assert.NoError(err, "unexpected error: %#v", err)
}

// GetOrderByID

// Negative tests

func TestGetOrderByID_InvalidID(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockOrderRepository(ctrl)
	logger := zerolog.Nop()
	h := handler.NewOrderHandler(repo, logger)

	id := ""

	r, _ := http.NewRequest(http.MethodGet, "/order/"+id, nil)
	w := httptest.NewRecorder()

	h.GetOrderByID(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	assert.Contains(w.Body.String(), errs.ErrInvalidOrderID.Error())
}

func TestGetOrderByID_OrderNotFound(t *testing.T) {
	t.Parallel()

	assert := a.New(t)

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockOrderRepository(ctrl)
	logger := zerolog.Nop()
	h := handler.NewOrderHandler(repo, logger)

	id := "123"

	r, _ := http.NewRequest(http.MethodGet, "/order/"+id, nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{"id": id})

	repo.EXPECT().GetOrderByID(gomock.Any(), id).Return(nil, errs.ErrOrderNotFound)

	h.GetOrderByID(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(http.StatusNotFound, resp.StatusCode)

	assert.Contains(w.Body.String(), errs.ErrOrderNotFound.Error())
	assert.Contains(w.Body.String(), errs.ErrGetOrderByID.Error())
}
