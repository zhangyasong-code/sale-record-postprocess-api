package models

import (
	"time"
)

type SaleRecordEvent struct {
	TransactionId             int64                   `json:"transactionId"`
	AssortedSaleRecordDtlList []AssortedSaleRecordDtl `json:"assortedSaleRecordDtlList"`
	TenantCode                string                  `json:"tenantCode"`
	StoreId                   int64                   `json:"storeId"`
	OrderId                   int64                   `json:"orderId"`
	OuterOrderNo              string                  `json:"outerOrderNo"`
	RefundId                  int64                   `json:"refundId"`
	TransactionType           string                  `json:"transactionType"`
	TransactionChannelType    string                  `json:"transactionChannelType"`
	TransactionStatus         string                  `json:"transactionStatus"`
	TransactionCreateDate     time.Time               `json:"transactionCreateDate"`
	TransactionUpdateDate     time.Time               `json:"transactionUpdateDate"`
	CustomerId                int64                   `json:"customerId"`
	SalesmanId                int64                   `json:"salesmanId"`
	TotalPrice                TotalPrice              `json:"totalPrice"`
	FreightPrice              float64                 `json:"freightPrice"`
	Mileage                   float64                 `json:"mileage"`
	MileagePrice              float64                 `json:"mileagePrice"`
	CashPrice                 float64                 `json:"cashPrice"`
	IsOutPaid                 bool                    `json:"isOutPaid"`
	CartOffers                []CartOffer             `json:"cartOffers"`
	Committed                 Committed               `json:"committed"`
	Payments                  []Payment               `json:"payments"`
}

type AssortedSaleRecordDtl struct {
	Id               int64            `json:"id"`
	OrderItemId      int64            `json:"orderItemId"`
	RefundItemId     int64            `json:"refundItemId"`
	BrandId          int64            `json:"brandId"`
	BrandCode        string           `json:"brandCode"`
	ItemCode         string           `json:"itemCode"`
	ItemName         string           `json:"itemName"`
	ProductId        int64            `json:"productId"`
	SkuId            int64            `json:"skuId"`
	SkuImg           string           `json:"skuImg"`
	ListPrice        float64          `json:"listPrice"`
	SalePrice        float64          `json:"salePrice"`
	Quantity         float64          `json:"quantity"`
	TotalPrice       TotalPrice       `json:"totalPrice"`
	DistributedPrice DistributedPrice `json:"distributedPrice"`
	Status           string           `json:"status"`
	ItemFee          float64          `json:"itemFee"`
	FeeRate          float64          `json:"feeRate"`
	ItemOffers       []Offer          `json:"itemOffers"`
	Committed        Committed        `json:"committed"`
}

type Committed struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"createdBy"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modifiedBy"`
}

type Offer struct {
	OfferId   int64   `json:"offerId"`
	OfferNo   string  `json:"offerNo"`
	CouponNo  string  `json:"couponNo"`
	ItemCodes string  `json:"itemCodes"`
	ItemCode  string  `json:"itemCode"`
	Price     float64 `json:"price"`
}

type TotalPrice struct {
	ListPrice        float64 `json:"listPrice"`
	SalePrice        float64 `json:"salePrice"`
	DiscountPrice    float64 `json:"discountPrice"`
	TransactionPrice float64 `json:"transactionPrice"`
}

type CartOffer struct {
	OfferId         int64   `json:"offerId"`
	OfferNo         string  `json:"offerNo"`
	CouponNo        string  `json:"couponNo"`
	ItemCodes       string  `json:"itemCodes"`
	TargetItemCodes string  `json:"targetItemCodes"`
	Price           float64 `json:"price"`
}

type Payment struct {
	SeqNo     int64     `json:"seqNo"`
	PayMethod string    `json:"payMethod"`
	PayAmt    float64   `json:"payAmt"`
	CreatedAt time.Time `json:"createdAt"`
}

type DistributedPrice struct {
	TotalDistributedItemOfferPrice float64 `json:"totalDistributedItemOfferPrice"`
	TotalDistributedCartOfferPrice float64 `json:"totalDistributedCartOfferPrice"`
	TotalDistributedPaymentPrice   float64 `json:"totalDistributedPaymentPrice"`
	DistributedCashPrice           float64 `json:"distributedCashPrice"`
}
