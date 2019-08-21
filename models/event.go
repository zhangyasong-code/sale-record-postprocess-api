package models

import "time"

type SaleRecordEvent struct {
	TransactionId            string                    `json:"transactionId"`
	TenantCode               string                    `json:"tenantCode"`
	StoreId                  int64                     `json:"storeId"`
	OrderId                  int64                     `json:"orderId"`
	OuterOrderNo             string                    `json:"outerOrderNo"`
	RefundId                 int64                     `json:"refundId"`
	TransactionType          string                    `json:"transactionType"`
	TransactionChannelType   string                    `json:"transactionChannelType"`
	TransactionStatus        string                    `json:"transactionStatus"`
	TransactionCreateDate    *time.Time                `json:"transactionCreateDate"`
	TransactionUpdateDate    string                    `json:"transactionUpdateDate"`
	CustomerId               int64                     `json:"customerId"`
	SalesmanId               int64                     `json:"salesmanId"`
	DiscountPrice            string                    `json:"discountPrice"`
	TotalPrice               *TotalPrice               `json:"totalPrice"`
	FreightPrice             string                    `json:"freightPrice"`
	Mileage                  float64                   `json:"mileage"`
	MileagePrice             float64                   `json:"mileagePrice"`
	CashPrice                string                    `json:"cashPrice"`
	IsOutPaid                bool                      `json:"isOutPaid"`
	IsRefund                 bool                      `json:"isRefund"`
	IsDelivery               bool                      `json:"isDelivery"`
	Committed                *Committed                `json:"committed"`
	AssortedSaleRecordDtls   []AssortedSaleRecordDtl   `json:"assortedSaleRecordDtls"`
	AssortedSaleRecordOffers []AssortedSaleRecordOffer `json:"assortedSaleRecordOffers"`
}

type AssortedSaleRecordDtl struct {
	Id                              string     `json:"id"`
	OrderItemId                     int64      `json:"orderItemId"`
	RefundItemId                    int64      `json:"refundItemId"`
	BrandId                         int64      `json:"brandId"`
	BrandCode                       string     `json:"brandCode"`
	ItemCode                        string     `json:"itemCode"`
	ItemName                        string     `json:"itemName"`
	ProductId                       string     `json:"productId"`
	SkuId                           int64      `json:"skuId"`
	SkuImg                          string     `json:"skuImg"`
	ListPrice                       float64    `json:"listPrice"`
	SalePrice                       float64    `json:"salePrice"`
	Quantity                        int        `json:"quantity"`
	TotalListPrice                  float64    `json:"totalListPrice"`
	TotalSalePrice                  float64    `json:"totalSalePrice"`
	TotalDistributedCartOfferPrice  float64    `json:"totalDistributedCartOfferPrice"`
	TotalDistributedCartCouponPrice float64    `json:"totalDistributedCartCouponPrice"`
	TotalDistributedPaymentPrice    float64    `json:"totalDistributedPaymentPrice"`
	DistributedCashPrice            float64    `json:"distributedCashPrice"`
	Status                          string     `json:"status"`
	Committed                       *Committed `json:"committed"`
}

type TotalPrice struct {
	ListPrice        float64 `json:"listPrice"`
	SalePrice        float64 `json:"salePrice"`
	DiscountPrice    float64 `json:"discountPrice"`
	TransactionPrice float64 `json:"transactionPrice"`
}

type Committed struct {
	Created    string `json:"created"`
	CreatedBy  string `json:"createdBy"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modifiedBy"`
}

type AssortedSaleRecordOffer struct{}
