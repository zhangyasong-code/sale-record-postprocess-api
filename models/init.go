package models

import "github.com/go-xorm/xorm"

func InitOrderDb(db *xorm.Engine) error {
	return nil
}

func InitSaleRecordDb(db *xorm.Engine) error {
	return db.Sync(
		new(PostSaleRecordFee),
		new(PostFailCreateSaleFee),
		new(PostMileage),
		new(PostMileageDtl),
		new(SaleRecordDtlSalesmanAmount),
		new(SaleRecordDtlOffer))
}
