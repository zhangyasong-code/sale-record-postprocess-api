package models

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
)

type PostSaleRecordFee struct {
	Id                     int64   `json:"id" query:"id"`
	TransactionId          string  `json:"transactionId" query:"transactionId" xorm:"index VARCHAR(30) notnull" validate:"required"`
	SaleRecordDtlId        int64   `json:"saleRecordDtlId" query:"saleRecordDtlId" xorm:"index VARCHAR(30) notnull" validate:"required"`
	SaleRecordOfferId      int64   `json:"saleRecordOfferId" query:"saleRecordOfferId" xorm:"index notnull" validate:"gte=0"`
	OrderId                int64   `json:"orderId" query:"orderId" xorm:"index notnull" validate:"gte=0"`
	OrderItemId            int64   `json:"orderItemId" query:"orderItemId" xorm:"index notnull" validate:"gte=0"`
	RefundId               int64   `json:"refundId" query:"refundId" xorm:"index notnull" validate:"gte=0"`
	RefundItemId           int64   `json:"refundItemId" query:"refundItemId" xorm:"index notnull" validate:"gte=0"`
	CustomerId             int64   `json:"customerId" query:"customerId" xorm:"index notnull" validate:"gte=0"`
	StoreId                int64   `json:"storeId" query:"storeId" xorm:"index notnull" validate:"gte=0"`
	EventType              string  `json:"eventType" query:"eventType" xorm:"VARCHAR(20)"`
	TotalSalePrice         float64 `json:"totalSalePrice" query:"totalSalePrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	TotalPaymentPrice      float64 `json:"totalPaymentPrice" query:"totalPaymentPrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	Mileage                float64 `json:"mileage" query:"mileage" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	MileagePrice           float64 `json:"mileagePrice" query:"mileagePrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	ContractFeeRate        float64 `json:"contractFeeRate" query:"contractFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	EventFeeRate           float64 `json:"eventFeeRate" query:"eventFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	AppliedFeeRate         float64 `json:"appliedFeeRate" query:"appliedFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	FeeAmount              float64 `json:"feeAmount" query:"feeAmount" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	TransactionChannelType string  `json:"transactionChannelType" query:"transactionChannelType" xorm:"VARCHAR(20)"`
}

func (p *PostSaleRecordFee) Save(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(p); err != nil {
		return err
	}
	return nil
}
