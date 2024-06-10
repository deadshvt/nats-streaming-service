package database

import (
	"context"

	"github.com/deadshvt/nats-streaming-service/internal/database/postgres"
	"github.com/deadshvt/nats-streaming-service/internal/entity"
	"github.com/deadshvt/nats-streaming-service/internal/errs"
)

type OrderDB interface {
	Connect() error
	Disconnect() error

	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrderByID(ctx context.Context, id string) (*entity.Order, error)
	GetAllOrders(ctx context.Context) ([]*entity.Order, error)
}

func InitOrderDB(dbType string) (OrderDB, error) {
	var db OrderDB

	switch dbType {
	case "postgres":
		db = &postgres.Postgres{}
	default:
		return nil, errs.ErrUnsupportedDBType
	}

	err := db.Connect()
	if err != nil {
		return nil, err
	}

	return db, nil
}
