package factory

import (
	"context"

	"github.com/go-xorm/xorm"
)

var (
	SaleRecordDBContextName = "saleRecordDB"
	OrderDBContextName      = "orderDB"
)

func SaleRecordDB(ctx context.Context) xorm.Interface {
	v := ctx.Value(SaleRecordDBContextName)
	// v := ctx.Value(echomiddleware.ContextDBName)
	if v == nil {
		panic("DB is not exist")
	}
	db, ok := v.(xorm.Interface)
	if !ok {
		panic("DB is not exist")
	}
	return db
}

func OrderDB(ctx context.Context) xorm.Interface {
	v := ctx.Value(OrderDBContextName)
	if v == nil {
		panic("DB is not exist")
	}
	db, ok := v.(xorm.Interface)
	if !ok {
		panic("DB is not exist")
	}
	return db
}
