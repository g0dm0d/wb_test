package nats

import (
	"github.com/g0dm0d/wbtest/internal/store"
	"github.com/g0dm0d/wbtest/pkg/cache"
	"github.com/nats-io/stan.go"
)

type Nats interface {
	HandleData(m *stan.Msg)
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
