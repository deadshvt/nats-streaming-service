package postgres

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/deadshvt/nats-streaming-service/internal/entity"
	"github.com/deadshvt/nats-streaming-service/internal/errs"

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

func (p *Postgres) CreateOrder(order *entity.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if tErr := tx.Rollback(); tErr != nil {
			p.Logger.Error().Msgf("Failed to rollback transaction: %v", tErr)
			if err == nil {
				err = tErr
			} else {
				err = errs.WrapError(err, tErr)
			}
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO "order" (data) VALUES ($1)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetOrderByID(id string) (*entity.Order, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if tErr := tx.Rollback(); tErr != nil {
			p.Logger.Error().Msgf("Failed to rollback transaction: %v", tErr)
			if err == nil {
				err = tErr
			} else {
				err = errs.WrapError(err, tErr)
			}
		}
	}()

	stmt, err := tx.Prepare(`
		SELECT data FROM "order" WHERE data->>('order_uid') = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var result []byte
	err = stmt.QueryRow(id).Scan(&result)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	var order entity.Order
	err = json.Unmarshal(result, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (p *Postgres) GetAllOrders() ([]*entity.Order, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if tErr := tx.Rollback(); tErr != nil {
			p.Logger.Error().Msgf("Failed to rollback transaction: %v", tErr)
			if err == nil {
				err = tErr
			} else {
				err = errs.WrapError(err, tErr)
			}
		}
	}()

	stmt, err := tx.Prepare(`
		SELECT data FROM "order"
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order

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
