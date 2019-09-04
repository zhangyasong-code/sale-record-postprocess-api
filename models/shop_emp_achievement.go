package models

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
)

type SaleRecordDtlSalesmanAmount struct {
	Id                          int64   `json:"id"`
	TransactionId               int64   `json:"transactionId"`
	SaleRecordDtlId             int64   `json:"saleRecordDtlId"`
	OrderId                     int64   `json:"orderId"`  //销售单号：退货时销售单号为原销售单号
	RefundId                    int64   `json:"refundId"` //退货单号：销售时退货单号为0
	StoreId                     int64   `json:"storeId"`
	SalesmanId                  int64   `json:"salesmanId"`
	TotalListPrice              float64 `json:"totalListPrice" xorm:"decimal(19,2)"`
	TotalSalePrice              float64 `json:"totalSalePrice" xorm:"decimal(19,2)"`
	TotalDiscountPrice          float64 `json:"totalDiscountPrice" xorm:"decimal(19,2)"`
	TotalDiscountCartOfferPrice float64 `json:"totalDiscountCartOfferPrice" xorm:"decimal(19,2)"`
	TotalPaymentPrice           float64 `json:"totalDistributedPaymentPrice" xorm:"decimal(19,2)"`
	TransactionChannelType      string  `json:"transactionChannelType"`
	Mileage                     float64 `json:"mileage" xorm:"decimal(19,2)"`      //mileagePrice.Point
	MileagePrice                float64 `json:"mileagePrice" xorm:"decimal(19,2)"` //mileagePrice.PointAmount
	SalesmanSaleAmount          float64 `json:"salesmanSaleAmount" xorm:"decimal(19,2)"`
	SalesmanSaleDiscountRate    float64 `json:"salesmanSaleDiscountRate" xorm:"decimal(19,2)"` //实际销售额活动扣率
	SalesmanNormalSaleAmount    float64 `json:"salesmanNormalSaleAmount" xorm:"decimal(19,2)"`
	SalesmanDiscountSaleAmount  float64 `json:"salesmanDiscountSaleAmount" xorm:"decimal(19,2)"`
}
type SaleRecordDtlOffer struct {
	Id               int64   `json:"id"`
	OfferId          int64   `json:"offerId"`
	SalesmanAmountId int64   `json:"salesmanAmountId"`
	EventType        string  `json:"eventType"`                            //PromotionEvent.EventTypeCode
	SaleBaseAmt      float64 `json:"saleBaseAmt" xorm:"decimal(19,2)"`     //PromotionEvent.SaleBaseAmt
	DiscountBaseAmt  float64 `json:"discountBaseAmt" xorm:"decimal(19,2)"` //PromotionEvent.DiscountBaseAmt
	DiscountRate     float64 `json:"discountRate" xorm:"decimal(19,2)"`    //PromotionEvent.DiscountRate
}

func (s *SaleRecordDtlSalesmanAmount) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(s); err != nil {
		return err
	}
	return nil
}
func (s *SaleRecordDtlOffer) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(s); err != nil {
		return err
	}
	return nil
}
