package models

import (
	"context"
)

type CustomerEventHandler struct {
}

func (h CustomerEventHandler) Handle(ctx context.Context, e SaleRecordEvent) error {

	// factory.OrderDB(ctx).Get(...)
	// factory.SaleRecordDB(ctx).Insert(...)

	return nil
}
