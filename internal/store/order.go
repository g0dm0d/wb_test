package store

import "github.com/g0dm0d/wbtest/internal/dto"

type SaveOrderOpts struct {
	Jsonb []byte
}

type GetOrderOpts struct {
	OrderID string
}

type OrderStore interface {
	SaveOrder(opts SaveOrderOpts) error
	GetOrder(opts GetOrderOpts) (dto.Order, error)
}
