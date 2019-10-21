package salePerson

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"time"
)

type SaleRecordDtlSalesmanAmount struct {
	Id                          int64     `json:"id"`
	TransactionId               int64     `json:"transactionId" xorm:"index"`
	SaleRecordDtlId             int64     `json:"saleRecordDtlId" xorm:"index"`
	OrderId                     int64     `json:"orderId"` //销售单号：退货时销售单号为原销售单号
	OrderItemId                 int64     `json:"orderItemId"`
	RefundId                    int64     `json:"refundId"` //退货单号：销售时退货单号为0
	RefundItemId                int64     `json:"refundItemId"`
	StoreId                     int64     `json:"storeId"`
	SalesmanId                  int64     `json:"salesmanId"`
	ItemCode                    string    `json:"itemCode"`
	TotalListPrice              float64   `json:"totalListPrice" xorm:"decimal(19,2)"`
	TotalSalePrice              float64   `json:"totalSalePrice" xorm:"decimal(19,2)"`
	TotalDiscountPrice          float64   `json:"totalDiscountPrice" xorm:"decimal(19,2)"`
	TotalDiscountItemOfferPrice float64   `json:"totalDiscountItemOfferPrice" xorm:"decimal(19,2)"`
	TotalDiscountCartOfferPrice float64   `json:"totalDiscountCartOfferPrice" xorm:"decimal(19,2)"`
	TotalPaymentPrice           float64   `json:"totalDistributedPaymentPrice" xorm:"decimal(19,2)"`
	TransactionType             string    `json:"transactionType"`
	TransactionChannelType      string    `json:"transactionChannelType"`
	Mileage                     float64   `json:"mileage" xorm:"decimal(19,2)"`      //mileagePrice.Point
	MileagePrice                float64   `json:"mileagePrice" xorm:"decimal(19,2)"` //mileagePrice.PointPrice
	SalesmanSaleAmount          float64   `json:"salesmanSaleAmount" xorm:"decimal(19,2)"`
	SalesmanSaleDiscountRate    float64   `json:"salesmanSaleDiscountRate" xorm:"decimal(19,2)"` //实际销售额活动扣率
	SalesmanNormalSaleAmount    float64   `json:"salesmanNormalSaleAmount" xorm:"decimal(19,2)"`
	SalesmanDiscountSaleAmount  float64   `json:"salesmanDiscountSaleAmount" xorm:"decimal(19,2)"`
	TransactionCreateDate       time.Time `json:"transactionCreateDate"` //销售时间
	CreatedAt                   time.Time `json:"createdAt" xorm:"created"`
}
type SaleRecordDtlOffer struct {
	Id               int64     `json:"id"`
	ItemCode         string    `json:"itemCode"`
	OfferId          int64     `json:"offerId" xorm:"index"`
	SalesmanAmountId int64     `json:"salesmanAmountId" xorm:"index"`
	EventType        string    `json:"eventType"`                            //PromotionEvent.EventTypeCode
	SaleBaseAmt      float64   `json:"saleBaseAmt" xorm:"decimal(19,2)"`     //PromotionEvent.SaleBaseAmt
	DiscountBaseAmt  float64   `json:"discountBaseAmt" xorm:"decimal(19,2)"` //PromotionEvent.DiscountBaseAmt
	DiscountRate     float64   `json:"discountRate" xorm:"decimal(19,2)"`    //PromotionEvent.DiscountRate
	CreatedAt        time.Time `json:"createdAt" xorm:"created"`
}

func (s *SaleRecordDtlSalesmanAmount) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(s); err != nil {
		return err
	}
	return nil
}
func (s *SaleRecordDtlSalesmanAmount) GetByKey(ctx context.Context, transactionId, saleRecordDtlId int64) (has bool, res *SaleRecordDtlSalesmanAmount, err error) {
	res = &SaleRecordDtlSalesmanAmount{}
	has, err = factory.SaleRecordDB(ctx).
		Where("transaction_id=?", transactionId).
		And("sale_record_dtl_id=?", saleRecordDtlId).
		Get(res)
	return
}
func (s *SaleRecordDtlOffer) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(s); err != nil {
		return err
	}
	return nil
}
func (s *SaleRecordDtlOffer) GetByKey(ctx context.Context, offerId, salesmanAmountId int64) (has bool, res *SaleRecordDtlOffer, err error) {
	res = &SaleRecordDtlOffer{}
	has, err = factory.SaleRecordDB(ctx).
		Where("offer_id=?", offerId).
		And("salesman_amount_id=?", salesmanAmountId).
		Get(res)
	return
}
func (e *SaleRecordDtlSalesmanAmount) Update(ctx context.Context, transactionId, saleRecordDtlId int64) (affectedRow int64, err error) {
	affectedRow, err = factory.SaleRecordDB(ctx).
		Where("transaction_id=?", transactionId).
		And("sale_record_dtl_id=?", saleRecordDtlId).Update(e)
	return
}
func (e *SaleRecordDtlOffer) Update(ctx context.Context, offerId, salesmanAmountId int64) (affectedRow int64, err error) {
	affectedRow, err = factory.SaleRecordDB(ctx).
		Where("offer_id=?", offerId).
		And("salesman_amount_id=?", salesmanAmountId).Update(e)
	return
}
