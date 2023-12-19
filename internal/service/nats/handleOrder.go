package nats

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/g0dm0d/wbtest/internal/dto"
	"github.com/g0dm0d/wbtest/internal/store"
	stan "github.com/nats-io/stan.go"
)

func (s *Service) HandleData(m *stan.Msg) {
	var message dto.Order

	if !CompareJSONToStruct(m.Data, message) {
		return
	}

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

	s.cache.Set(message.OrderUid, m.Data)
	return
}

func CompareJSONToStruct(bytes []byte, empty interface{}) bool {
	var mapped map[string]interface{}

	if err := json.Unmarshal(bytes, &mapped); err != nil {
		return false
	}

	emptyValue := reflect.ValueOf(empty).Type()

	if len(mapped) != emptyValue.NumField() {
		return false
	}

	for key := range mapped {
		if field, found := emptyValue.FieldByName(key); found {
			if !strings.EqualFold(key, strings.Split(field.Tag.Get("json"), ",")[0]) {
				return false
			}
		}
	}

	return true
}
