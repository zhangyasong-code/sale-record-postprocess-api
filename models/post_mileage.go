package models

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"time"
)

type UseType string

const (
	UseTypeEarn       UseType = "Earn"
	UseTypeEarnCancel UseType = "EarnCancel"
	UseTypeUsed       UseType = "Used"
	UseTypeUsedCancel UseType = "UsedCancel"
)

type PostMileage struct {
	Id            int64     `json:"id"`
	TenantCode    string    `json:"tenantCode" xorm:"index VARCHAR(50) notnull" validate:"required"`
	StoreId       int64     `json:"storeId" xorm:"index notnull" validate:"gte=0"`
	CustomerId    int64     `json:"customerId" xorm:"index notnull" validate:"gte=0"`
	BrandId       int64     `json:"brandId" xorm:"notnull"`
	GradeId       int64     `json:"gradeId"`
	TransactionId int64     `json:"transactionId" xorm:"index default 0" validate:"required"`
	OrderId       int64     `json:"orderId" xorm:"index default 0" validate:"required"`
	RefundId      int64     `json:"refundId" xorm:"index default 0"`
	UseType       UseType   `json:"useType" xorm:"VARCHAR(25)"`
	Point         float64   `json:"point" xorm:"decimal(19,2)"`
	PointAmount   float64   `json:"pointAmount" xorm:"decimal(19,2)"`
	CreateAt      time.Time `json:"createAt" xorm:"created"`
	UpdateAt      time.Time `json:"updateAt" xorm:"updated"`
}

type PostMileageDtl struct {
	Id               int64      `json:"id"`
	PostMileageId    int64      `json:"postMileageId"`
	TransactionDtlId int64      `json:"transactionDtlId" xorm:"index default 0"`
	OrderItemId      int64      `json:"orderItemId" xorm:"index default 0"`
	RefundItemId     int64      `json:"refundItemId" xorm:"index default 0"`
	UseType          UseType    `json:"useType" xorm:"VARCHAR(25)"`
	Point            float64    `json:"point" xorm:"decimal(19,2)"`
	PointAmount      float64    `json:"pointAmount" xorm:"decimal(19,2)"`
	CreateAt         *time.Time `json:"createAt" xorm:"created"`
	UpdateAt         *time.Time `json:"updateAt" xorm:"updated"`
}

func (PostMileage) MakePostMileage(mileage *Mileage, mileageMall *MileageMall, record SaleRecordEvent) *PostMileage {
	var useType UseType
	point := record.Mileage
	pointAmount := record.MileagePrice
	if mileage == nil {
		point = mileage.Point
		pointAmount = 0
		if record.RefundId != 0 {
			useType = UseTypeUsedCancel
		} else {
			useType = UseTypeUsed
		}
	} else {
		if record.RefundId != 0 {
			useType = UseTypeEarnCancel
		} else {
			useType = UseTypeEarn
		}
	}
	postMileage := &PostMileage{
		TenantCode:    record.TenantCode,
		StoreId:       record.StoreId,
		CustomerId:    record.CustomerId,
		BrandId:       mileageMall.MallId,
		GradeId:       mileageMall.GradeId,
		TransactionId: record.TransactionId,
		OrderId:       record.OrderId,
		RefundId:      record.RefundId,
		UseType:       useType,
		Point:         point,
		PointAmount:   pointAmount,
	}
	return postMileage
}

func (o *PostMileage) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(o); err != nil {
		return err
	}
	return nil
}

func (PostMileage) CheckOrderRefundExist(ctx context.Context, transactionId int64) (bool, error) {
	postMileage := PostMileage{}
	has, err := factory.SaleRecordDB(ctx).Table("post_mileage").
		Where("transaction_id=?", transactionId).Get(&postMileage)
	if err != nil {
		return true, err
	}
	return has, nil
}

func (PostMileageDtl) MakePostMileageDtl(postMileageId int64, useType UseType, point, pointAmount float64, recordDtl AssortedSaleRecordDtl) *PostMileageDtl {
	postMileageDtl := &PostMileageDtl{
		PostMileageId:    postMileageId,
		TransactionDtlId: recordDtl.Id,
		OrderItemId:      recordDtl.OrderItemId,
		RefundItemId:     recordDtl.RefundItemId,
		UseType:          useType,
		Point:            point,
		PointAmount:      pointAmount,
	}
	return postMileageDtl
}
func (PostMileageDtl) CreateBatch(ctx context.Context, v []PostMileageDtl) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(v); err != nil {
		return err
	}
	return nil
}
func (PostMileageDtl) GetByKey(ctx context.Context, transactionDtlId, orderItemId, refundItemId int64) (has bool, res *PostMileageDtl, err error) {
	res = &PostMileageDtl{}
	has, err = factory.SaleRecordDB(ctx).
		Where("transaction_dtl_id=?", transactionDtlId).
		And("order_item_id=?", orderItemId).
		And("refund_item_id=?", refundItemId).
		Get(res)
	return
}
