package order

import (
	"github.com/g0dm0d/wbtest/internal/server/req"
	"github.com/g0dm0d/wbtest/internal/store"
	"github.com/g0dm0d/wbtest/pkg/cache"
)

type Order interface {
	GetOrder(ctx *req.Ctx) error
}

type Service struct {
	orderStore store.OrderStore
	cache      *cache.Map
}

func New(orderStore store.OrderStore, cacheMap *cache.Map) *Service {
	return &Service{
		orderStore: orderStore,
		cache:      cacheMap,
	}
}
