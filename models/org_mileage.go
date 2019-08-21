package models

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"time"
)

type MileageType string

const (
	MileageTypeEarn       MileageType = "Earn"
	MileageTypeEarnCancel MileageType = "EarnCancel"
	MileageTypeUsed       MileageType = "Used"
	MileageTypeUsedCancel MileageType = "UsedCancel"
)

type OrgMileage struct {
	Id          int64       `json:"id"`
	TenantCode  string      `json:"tenantCode" xorm:"index VARCHAR(50) notnull" validate:"required"`
	StoreId     int64       `json:"storeId" xorm:"index notnull" validate:"gte=0"`
	CustomerId  int64       `json:"customerId" xorm:"index notnull" validate:"gte=0"`
	OrderId     int64       `json:"orderId" xorm:"index default 0" validate:"required"`
	RefundId    int64       `json:"refundId" xorm:"index default 0"`
	MileageType MileageType `json:"mileageType" xorm:"VARCHAR(25)"`
	Point       float64     `json:"point" xorm:"decimal(19,2)"`
	PointAmount float64     `json:"pointAmount" xorm:"decimal(19,2)"`
	CreateAt    time.Time   `json:"createAt" xorm:"created"`
	UpdateAt    time.Time   `json:"updateAt" xorm:"updated"`
}

type OrgMileageDtl struct {
	Id           int64       `json:"id"`
	MstId        int64       `json:"mstId"`
	OrderItemId  int64       `json:"orderItemId" xorm:"index default 0"`
	RefundItemId int64       `json:"refundItemId" xorm:"index default 0"`
	MileageType  MileageType `json:"mileageType" xorm:"VARCHAR(25)"`
	Point        float64     `json:"point" xorm:"decimal(19,2)"`
	PointAmount  float64     `json:"pointAmount" xorm:"decimal(19,2)"`
	CreateAt     *time.Time  `json:"createAt" xorm:"created"`
	UpdateAt     *time.Time  `json:"updateAt" xorm:"updated"`
}

func (OrgMileage) MakeOrgMileage(mileage *Mileage, record SaleRecordEvent) *OrgMileage {
	var mileageType MileageType
	point := record.Mileage
	pointAmount := record.MileagePrice
	if mileage == nil {
		point = mileage.Point
		pointAmount = 0
		if record.RefundId != 0 {
			mileageType = MileageTypeUsedCancel
		} else {
			mileageType = MileageTypeUsed
		}
	} else {
		if record.RefundId != 0 {
			mileageType = MileageTypeEarnCancel
		} else {
			mileageType = MileageTypeEarn
		}
	}
	orgMileage := &OrgMileage{
		TenantCode:  record.TenantCode,
		StoreId:     record.StoreId,
		CustomerId:  record.CustomerId,
		OrderId:     record.OrderId,
		RefundId:    record.RefundId,
		MileageType: mileageType,
		Point:       point,
		PointAmount: pointAmount,
	}
	return orgMileage
}

func (o *OrgMileage) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(o); err != nil {
		return err
	}
	return nil
}

func (OrgMileageDtl) MakeOrgMileageDtl(mstId int64, mileageType MileageType, point, pointAmount float64, recordDtl AssortedSaleRecordDtl) *OrgMileageDtl {
	orgMileageDtl := &OrgMileageDtl{
		MstId:        mstId,
		OrderItemId:  recordDtl.OrderItemId,
		RefundItemId: recordDtl.RefundItemId,
		MileageType:  mileageType,
		Point:        point,
		PointAmount:  pointAmount,
	}
	return orgMileageDtl
}
func (OrgMileageDtl) CreateBatch(ctx context.Context, v []OrgMileageDtl) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(v); err != nil {
		return err
	}
	return nil
}
