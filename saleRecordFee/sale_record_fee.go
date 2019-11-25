package saleRecordFee

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"time"
)

type PostSaleRecordFee struct {
	TransactionDtlId       int64     `json:"transactionDtlId" query:"transactionDtlId" xorm:"index notnull pk" validate:"required"`
	TransactionId          int64     `json:"transactionId" query:"transactionId" xorm:"index notnull" validate:"required"`
	OrderId                int64     `json:"orderId" query:"orderId" xorm:"index notnull" validate:"gte=0"`
	OrderItemId            int64     `json:"orderItemId" query:"orderItemId" xorm:"index notnull" validate:"gte=0"`
	RefundId               int64     `json:"refundId" query:"refundId" xorm:"index notnull" validate:"gte=0"`
	RefundItemId           int64     `json:"refundItemId" query:"refundItemId" xorm:"index notnull" validate:"gte=0"`
	CustomerId             int64     `json:"customerId" query:"customerId" xorm:"index notnull" validate:"gte=0"`
	StoreId                int64     `json:"storeId" query:"storeId" xorm:"index notnull" validate:"gte=0"`
	TotalSalePrice         float64   `json:"totalSalePrice" query:"totalSalePrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	TotalPaymentPrice      float64   `json:"totalPaymentPrice" query:"totalPaymentPrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	Mileage                float64   `json:"mileage" query:"mileage" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	MileagePrice           float64   `json:"mileagePrice" query:"mileagePrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	ItemFeeRate            float64   `json:"itemFeeRate" query:"itemFeeRate" xorm:"DECIMAL(18,4) default 0" validate:"gte=0"`
	ItemFee                float64   `json:"itemFee" query:"itemFee" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	EventFeeRate           float64   `json:"eventFeeRate" query:"eventFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	AppliedFeeRate         float64   `json:"appliedFeeRate" query:"appliedFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	FeeAmount              float64   `json:"feeAmount" query:"feeAmount" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	TransactionChannelType string    `json:"transactionChannelType" query:"transactionChannelType" xorm:"index VARCHAR(20) notnull"`
	CreatedAt              time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt              time.Time `json:"updatedAt" xorm:"updated"`
}

type PostFailCreateSaleFee struct {
	TransactionId int64     `json:"transactionId" query:"transactionId" xorm:"index notnull pk" validate:"required"`
	IsProcessed   bool      `json:"isProcessed" query:"isProcessed" xorm:"index notnull default false" validate:"required"`
	CreatedAt     time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt     time.Time `json:"updatedAt" xorm:"updated"`
}

type ResultShopMessage struct {
	Success bool `json:"success"`
	Result  struct {
		Items []struct {
			Contracts []Contract `json:"contracts"`
		} `json:"items"`
	} `json:"result"`
	Error struct{} `json:"error"`
}

type Contract struct {
	Id         int64     `json:"id"`
	StoreId    int64     `json:"storeId"`
	BrandId    int64     `json:"brandId"`
	ContractNo string    `json:"contractNo"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	BaseRate   float64   `json:"baseRate"`
}

func (p *PostSaleRecordFee) Save(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(p); err != nil {
		return err
	}
	return nil
}

func (p *PostSaleRecordFee) Get(ctx context.Context) (bool, PostSaleRecordFee, error) {
	postSaleRecordFee := PostSaleRecordFee{}
	has, err := factory.SaleRecordDB(ctx).Where("transaction_dtl_id = ?", p.TransactionDtlId).Get(&postSaleRecordFee)
	if err != nil {
		return has, PostSaleRecordFee{}, err
	}
	return has, postSaleRecordFee, nil
}

func (pf *PostFailCreateSaleFee) Save(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(pf); err != nil {
		return err
	}
	return nil
}

func (pf *PostFailCreateSaleFee) Get(ctx context.Context) (bool, PostFailCreateSaleFee, error) {
	postFailCreateSaleFee := PostFailCreateSaleFee{}
	has, err := factory.SaleRecordDB(ctx).Where("transaction_id = ?", pf.TransactionId).And("is_processed = ?", false).Get(&postFailCreateSaleFee)
	if err != nil {
		return has, PostFailCreateSaleFee{}, err
	}
	return has, postFailCreateSaleFee, nil
}
