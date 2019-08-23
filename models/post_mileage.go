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

type PostMileage struct {
	Id              int64       `json:"id"`
	TenantCode      string      `json:"tenantCode" xorm:"index VARCHAR(50) notnull" validate:"required"`
	StoreId         int64       `json:"storeId" xorm:"index notnull" validate:"gte=0"`
	CustomerId      int64       `json:"customerId" xorm:"index notnull" validate:"gte=0"`
	SaleRecordMstId string      `json:"saleRecordMstId" xorm:"index default 0" validate:"required"`
	OrderId         int64       `json:"orderId" xorm:"index default 0" validate:"required"`
	RefundId        int64       `json:"refundId" xorm:"index default 0"`
	MileageType     MileageType `json:"mileageType" xorm:"VARCHAR(25)"`
	Point           float64     `json:"point" xorm:"decimal(19,2)"`
	PointAmount     float64     `json:"pointAmount" xorm:"decimal(19,2)"`
	CreateAt        time.Time   `json:"createAt" xorm:"created"`
	UpdateAt        time.Time   `json:"updateAt" xorm:"updated"`
}

type PostMileageDtl struct {
	Id              int64       `json:"id"`
	MstId           int64       `json:"mstId"`
	SaleRecordDtlId int64       `json:"saleRecordDtlId" xorm:"index default 0"`
	OrderItemId     int64       `json:"orderItemId" xorm:"index default 0"`
	RefundItemId    int64       `json:"refundItemId" xorm:"index default 0"`
	MileageType     MileageType `json:"mileageType" xorm:"VARCHAR(25)"`
	Point           float64     `json:"point" xorm:"decimal(19,2)"`
	PointAmount     float64     `json:"pointAmount" xorm:"decimal(19,2)"`
	CreateAt        *time.Time  `json:"createAt" xorm:"created"`
	UpdateAt        *time.Time  `json:"updateAt" xorm:"updated"`
}

func (PostMileage) MakePostMileage(mileage *Mileage, record SaleRecordEvent) *PostMileage {
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
	postMileage := &PostMileage{
		TenantCode:      record.TenantCode,
		StoreId:         record.StoreId,
		CustomerId:      record.CustomerId,
		SaleRecordMstId: record.TransactionId,
		OrderId:         record.OrderId,
		RefundId:        record.RefundId,
		MileageType:     mileageType,
		Point:           point,
		PointAmount:     pointAmount,
	}
	return postMileage
}

func (o *PostMileage) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(o); err != nil {
		return err
	}
	return nil
}

func (PostMileageDtl) MakePostMileageDtl(mstId int64, mileageType MileageType, point, pointAmount float64, recordDtl AssortedSaleRecordDtl) *PostMileageDtl {
	postMileageDtl := &PostMileageDtl{
		MstId:           mstId,
		SaleRecordDtlId: recordDtl.Id,
		OrderItemId:     recordDtl.OrderItemId,
		RefundItemId:    recordDtl.RefundItemId,
		MileageType:     mileageType,
		Point:           point,
		PointAmount:     pointAmount,
	}
	return postMileageDtl
}
func (PostMileageDtl) CreateBatch(ctx context.Context, v []PostMileageDtl) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(v); err != nil {
		return err
	}
	return nil
}
