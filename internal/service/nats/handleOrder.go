package nats

import (
	"encoding/json"
	"log"

	"github.com/g0dm0d/wbtest/internal/dto"
	"github.com/g0dm0d/wbtest/internal/store"
	stan "github.com/nats-io/stan.go"
)

func (s *Service) HandleData(m *stan.Msg) {
	var message dto.Order
	err := json.Unmarshal(m.Data, &message)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.orderStore.SaveOrder(store.SaveOrderOpts{Jsonb: m.Data})
	if err != nil {
		log.Println(err)
		return
	}

	s.cache.Set(message.OrderUid, message)
	return
}
