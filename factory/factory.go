package factory

import (
	"context"

	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

var (
	SaleRecordDBContextName    echomiddleware.ContextDBType = "saleRecordDB"
	OrderDBContextName         echomiddleware.ContextDBType = "orderDB"
	Mslv2ReadonlyDBContextName echomiddleware.ContextDBType = "mslv2ReadonlyDB"
)

func SaleRecordDB(ctx context.Context) xorm.Interface {
	v := ctx.Value(SaleRecordDBContextName)
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

func Mslv2ReadonlyDB(ctx context.Context) xorm.Interface {
	v := ctx.Value(Mslv2ReadonlyDBContextName)
	if v == nil {
		panic("DB is not exist")
	}
	db, ok := v.(xorm.Interface)
	if !ok {
		panic("DB is not exist")
	}
	return db
}
