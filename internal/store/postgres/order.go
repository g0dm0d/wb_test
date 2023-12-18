package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/g0dm0d/wbtest/internal/dto"
	"github.com/g0dm0d/wbtest/internal/store"
)

type OrderStore struct {
	db *sql.DB
}

func NewOrderStore(db *sql.DB) store.OrderStore {
	return &OrderStore{
		db: db,
	}
}

func (s *OrderStore) SaveOrder(opts store.SaveOrderOpts) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO Orders(data) VALUES($1)`, opts.Jsonb)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderStore) GetOrder(opts store.GetOrderOpts) (dto.Order, error) {
	var order dto.Order
	var jsonb []byte

	row := s.db.QueryRow(`SELECT * FROM Orders
		WHERE data ->>'order_uid' = $1`, opts.OrderID)

	err := row.Scan(&jsonb)

	json.Unmarshal(jsonb, &order)

	if err != nil {
		return order, err
	}

	return order, nil
}
