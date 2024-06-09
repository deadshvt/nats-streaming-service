package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deadshvt/nats-streaming-service/internal/entity"
	"github.com/deadshvt/nats-streaming-service/internal/errs"
	"github.com/deadshvt/nats-streaming-service/internal/html"
	"github.com/deadshvt/nats-streaming-service/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

const (
	Pattern       = "./assets/html/*.html"
	OrderFileName = "order.html"
	IndexFileName = "index.html"
	ErrorFileName = "error.html"
)

type OrderHandler struct {
	Repository entity.OrderRepository
	Logger     zerolog.Logger
}

func NewOrderHandler(repository entity.OrderRepository, logger zerolog.Logger) *OrderHandler {
	return &OrderHandler{
		Repository: repository,
		Logger:     logger,
	}
}

func (h *OrderHandler) CreateOrder(data []byte) error {
	h.Logger.Info().Msg("Creating order...")

	var order entity.Order
	if err := json.Unmarshal(data, &order); err != nil {
		msg := errs.WrapError(errs.ErrJSONUnmarshal, err)
		h.Logger.Error().Msg(msg.Error())
		return msg
	}

	err := ValidateOrder(&order)
	if err != nil {
		msg := errs.WrapError(errs.ErrValidateOrder, err)
		h.Logger.Error().Msg(msg.Error())
		return msg
	}

	err = h.Repository.CreateOrder(&order)
	if err != nil {
		msg := errs.WrapError(errs.ErrCreateOrder, err)
		h.Logger.Error().Msg(msg.Error())
		return msg
	}

	logger.LogWithParams(h.Logger, "Created order", struct {
		OrderID string
	}{OrderID: order.OrderUid})

	return nil
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("Getting order by id...")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		msg := errs.WrapError(errs.ErrInvalidOrderID, nil).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	order, err := h.Repository.GetOrderByID(id)
	if err != nil {
		msg := errs.WrapError(errs.ErrGetOrderByID, err).Error()
		h.Logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	logger.LogWithParams(h.Logger, "Got order", struct {
		OrderID string
	}{OrderID: id})

	SetResponse(w, http.StatusOK)
	html.ParseTemplate(h.Logger, w, Pattern, OrderFileName, order)
}

func (h *OrderHandler) GetOrderID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("Getting order id...")

	if r.Method == http.MethodPost {
		orderID := r.FormValue("orderID")
		logger.LogWithParams(h.Logger, "Got orderID", struct {
			OrderID string
		}{OrderID: orderID})

		_, err := h.Repository.GetOrderByID(orderID)
		if err != nil {
			msg := errs.WrapError(errs.ErrGetOrderByID, err).Error()
			h.Logger.Error().Msg(msg)
			html.ParseTemplate(h.Logger, w, Pattern, ErrorFileName, nil)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/order/%s", orderID), http.StatusSeeOther)
		}
	} else {
		html.ParseTemplate(h.Logger, w, Pattern, IndexFileName, nil)
	}
}

func SetResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
}

func ValidateOrder(order *entity.Order) error {
	v := validator.New()
	if err := v.Struct(order); err != nil {
		return err
	}

	return nil
}
