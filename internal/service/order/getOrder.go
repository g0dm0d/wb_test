package order

import (
	"encoding/json"

	"github.com/g0dm0d/wbtest/internal/server/req"
	"github.com/g0dm0d/wbtest/internal/store"
	"github.com/g0dm0d/wbtest/pkg/errs"
	"github.com/go-chi/chi/v5"
)

func (s *Service) GetOrder(ctx *req.Ctx) (err error) {
	odrderID := chi.URLParam(ctx.Request, "orderID")

	value, ok := s.cache.Get(odrderID)
	if !ok {
		valueDB, err := s.orderStore.GetOrder(store.GetOrderOpts{
			OrderID: odrderID,
		})
		if err != nil {
			return errs.ReturnError(ctx.Writer, errs.InvalidID)
		}
		value, err = json.Marshal(valueDB)
		if err != nil {
			return errs.ReturnError(ctx.Writer, errs.InternalServerError)
		}
		s.cache.Set(odrderID, value)
	}

	ctx.Writer.Header().Set("Content-Type", "application/json")
	_, err = ctx.Writer.Write(value.([]byte))
	return err
}
