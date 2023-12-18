package order

import (
	"github.com/g0dm0d/wbtest/internal/dto"
	"github.com/g0dm0d/wbtest/internal/server/req"
	"github.com/g0dm0d/wbtest/internal/store"
	"github.com/g0dm0d/wbtest/pkg/errs"
	"github.com/go-chi/chi/v5"
)

func (s *Service) GetOrder(ctx *req.Ctx) (err error) {
	odrderID := chi.URLParam(ctx.Request, "orderID")

	value, ok := s.cache.Get(odrderID)
	if !ok {
		value, err = s.orderStore.GetOrder(store.GetOrderOpts{
			OrderID: odrderID,
		})
		if err != nil {
			return errs.ReturnError(ctx.Writer, errs.InvalidID)
		}
		s.cache.Set(odrderID, value)
	}
	return ctx.JSON(value.(dto.Order))
}
