package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/deadshvt/nats-streaming-service/internal/entity"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Postgres struct {
	DB     *sql.DB
	Logger zerolog.Logger
}

func (p *Postgres) Connect() error {
	p.Logger.Info().Msg("Connecting to Postgres...")

	dsn := os.Getenv("DB_DSN")
	var err error
	p.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil
	}

	return p.DB.Ping()
}

func (p *Postgres) Disconnect() error {
	p.Logger.Info().Msg("Disconnecting from Postgres...")

	if p.DB == nil {
		return nil
	}

	return p.DB.Close()
}

func (p *Postgres) CreateOrder(ctx context.Context, order *entity.Order) (err error) {
	p.Logger.Info().Msg("Creating order...")

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				p.Logger.Error().Msgf("Failed to rollback transaction: %v", rbErr)
				err = fmt.Errorf("%v and %v", err, rbErr)
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO "order" (data) VALUES ($1)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, data)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetOrderByID(ctx context.Context, id string) (order *entity.Order, err error) {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				p.Logger.Error().Msgf("Failed to rollback transaction: %v", rbErr)
				err = fmt.Errorf("%v and %v", err, rbErr)
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		SELECT data FROM "order" WHERE data->>('order_uid') = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var result []byte
	err = stmt.QueryRowContext(ctx, id).Scan(&result)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(result, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (p *Postgres) GetAllOrders(ctx context.Context) (orders []*entity.Order, err error) {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				p.Logger.Error().Msgf("Failed to rollback transaction: %v", rbErr)
				err = fmt.Errorf("%v and %v", err, rbErr)
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		SELECT data FROM "order"
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result []byte
		if err = rows.Scan(&result); err != nil {
			return nil, err
		}

		var order entity.Order
		if err = json.Unmarshal(result, &order); err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
