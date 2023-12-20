package service

import (
	"github.com/g0dm0d/wbtest/internal/service/nats"
	"github.com/g0dm0d/wbtest/internal/service/order"
	"github.com/g0dm0d/wbtest/internal/store"
	"github.com/g0dm0d/wbtest/pkg/cache"
)

type Service struct {
	Order order.Order
	Nats  nats.Nats
}

func New(orderStore store.OrderStore, cache *cache.Map) *Service {
	return &Service{
		Order: order.New(orderStore, cache),
		Nats:  nats.New(orderStore, cache),
	}
}
