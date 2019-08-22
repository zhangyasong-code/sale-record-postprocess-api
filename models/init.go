package models

import "github.com/go-xorm/xorm"

func InitOrderDb(db *xorm.Engine) error {
	return db.Sync(new(PostSaleRecordFee))
}

func InitSaleRecordDb(db *xorm.Engine) error {
	return db.Sync(
		new(OrgMileage),
		new(OrgMileageDtl))
}
