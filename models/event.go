package models

import "time"

type SaleRecordEvent struct {
	AuthToken string    `json:"authToken"`
	Payload   EventBody `json:"payload"`
}
type EventBody struct {
	TransactionId          string                  `json:"transactionId"`
	AssortedSaleRecordDtls []AssortedSaleRecordDtl `json:"assortedSaleRecordDtls"`
	TenantCode             string                  `json:"tenantCode"`
	StoreId                int64                   `json:"storeId"`
	OrderId                int64                   `json:"orderId"`
	OuterOrderNo           string                  `json:"outerOrderNo"`
	RefundId               int64                   `json:"refundId"`
	TransactionType        string                  `json:"transactionType"`
	TransactionChannelType string                  `json:"transactionChannelType"`
	TransactionStatus      string                  `json:"transactionStatus"`
	TransactionCreateDate  time.Time               `json:"transactionCreateDate"`
	TransactionUpdateDate  time.Time               `json:"transactionUpdateDate"`
	CustomerId             int64                   `json:"customerId"`
	SalesmanId             int64                   `json:"salesmanId"`
	DiscountPrice          float64                 `json:"discountPrice"`
	TotalPrice             TotalPrice              `json:"totalPrice"`
	FreightPrice           float64                 `json:"freightPrice"`
	Mileage                float64                 `json:"mileage"`
	MileagePrice           float64                 `json:"mileagePrice"`
	CashPrice              float64                 `json:"cashPrice"`
	IsOutPaid              bool                    `json:"isOutPaid"`
	IsRefund               bool                    `json:"isRefund"`
	IsDelivery             bool                    `json:"isDelivery"`
	Committed              Committed               `json:"committed"`
}

type AssortedSaleRecordDtl struct {
	Id                              int64     `json:"id"`
	OrderItemId                     int64     `json:"orderItemId"`
	RefundItemId                    int64     `json:"refundItemId"`
	BrandId                         int64     `json:"brandId"`
	BrandCode                       string    `json:"brandCode"`
	ItemCode                        string    `json:"itemCode"`
	ItemName                        string    `json:"itemName"`
	ProductId                       int64     `json:"productId"`
	SkuId                           int64     `json:"skuId"`
	SkuImg                          string    `json:"skuImg"`
	ListPrice                       float64   `json:"listPrice"`
	SalePrice                       float64   `json:"salePrice"`
	Quantity                        float64   `json:"quantity"`
	TotalListPrice                  float64   `json:"totalListPrice"`
	TotalSalePrice                  float64   `json:"totalSalePrice"`
	TotalDistributedCartOfferPrice  float64   `json:"totalDistributedCartOfferPrice"`
	TotalDistributedCartCouponPrice float64   `json:"totalDistributedCartCouponPrice"`
	TotalDistributedPaymentPrice    float64   `json:"totalDistributedPaymentPrice"`
	DistributedCashPrice            float64   `json:"distributedCashPrice"`
	Status                          string    `json:"status"`
	Offer                           Offer     `json:"offer"`
	Committed                       Committed `json:"committed"`
}

type Committed struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"createdBy"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modifiedBy"`
}

type Offer struct {
	Id             int64     `json:"id"`
	TransactionId  string    `json:"transactionId"`
	OrderId        int64     `json:"orderId"`
	OrderItemId    int64     `json:"orderItemId"`
	RefundId       int64     `json:"refundId"`
	RefundItemId   int64     `json:"refundItemId"`
	StoreId        int64     `json:"storeId"`
	CustomerId     int64     `json:"customerId"`
	SalesmanId     int64     `json:"salesmanId"`
	OfferId        int64     `json:"offerId"`
	OfferNo        string    `json:"offerNo"`
	CouponNo       string    `json:"couponNo"`
	ItemCodes      string    `json:"itemCodes"`
	Description    string    `json:"description"`
	OfferPrice     float64   `json:"offerPrice"`
	OfferCreatedAt time.Time `json:"offerCreatedAt"`
}

type TotalPrice struct {
	ListPrice        float64 `json:"listPrice"`
	SalePrice        float64 `json:"salePrice"`
	DiscountPrice    float64 `json:"discountPrice"`
	TransactionPrice float64 `json:"transactionPrice"`
}
