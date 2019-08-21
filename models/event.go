package models

import (
	"context"
	"errors"
	"strconv"
	"time"
)

type SaleRecordEvent struct {
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

type AssortedSaleRecordOffer struct{}
type Brand struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Store struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Promotion struct {
	EventNo         string  `json:"eventNo"`
	EventType       string  `json:"eventType"`
	SaleBaseAmt     float64 `json:"saleBaseAmt"`
	DiscountBaseAmt float64 `json:"discountBaseAmt"`
	DiscountRate    float64 `json:"discountRate"`
}

func getBrandAndStore(ctx context.Context, channels []ChannelCondition) (*Brand, *Store, error) {
	var (
		brandId, storeId int64
		err              error
	)
	for i := range channels {
		if channels[i].Type == "brand_id" {
			brandId, err = strconv.ParseInt(channels[i].Value, 10, 64)
			if err != nil {
				return nil, nil, errors.New("brandId invalid")
			}
		}
		if channels[i].Type == "store_id" {
			storeId, err = strconv.ParseInt(channels[i].Value, 10, 64)
			if err != nil {
				return nil, nil, errors.New("storeId invalid")
			}
		}
	}
	brand, err := getBrand(ctx, brandId)
	if err != nil {
		return nil, nil, err
	}
	if storeId == 0 {
		return brand, nil, err
	}
	store, err := getStore(ctx, storeId)
	if err != nil {
		return nil, nil, err
	}
	return brand, store, nil
}

func ToCSLOfferType(offerType OfferType, templateCode string) (string, error) {
	var eventTypeCode string
	switch templateCode {
	case "A1", "A5", "A7":
		eventTypeCode = "03"
	case "A2":
		if offerType == OfferTypeBrand {
			eventTypeCode = "03"
		} else if offerType == OfferTypeMember {
			eventTypeCode = "B"
		} else if offerType == OfferTypeChannel {
			eventTypeCode = "V"
		}
	case "B1", "B2", "C1", "C2", "C6", "D1":
		eventTypeCode = "02"
	case "C5":
		if offerType == OfferTypeBrand {
			eventTypeCode = "02"
		} else if offerType == OfferTypeMember {
			eventTypeCode = "P"
		}
	case "D2":
		if offerType == OfferTypeBrand {
			eventTypeCode = "02"
		} else if offerType == OfferTypeMember {
			eventTypeCode = "G"
		}
	case "D5":
		if offerType == OfferTypeBrand {
			eventTypeCode = "02"
		} else if offerType == OfferTypeMember {
			eventTypeCode = "R"
		}
	case "D6":
		eventTypeCode = "01"
	case "D7":
		eventTypeCode = "M"
	default:
		return eventTypeCode, errors.New("invalid offerType")
	}
	return eventTypeCode, nil
}

func (p *Promotion) ToCSLDisCount(StandardValue, DiscountValue float64) {
	switch p.EventType {
	case "01":
		p.SaleBaseAmt = StandardValue
		p.DiscountBaseAmt = DiscountValue
		p.EventNo = "12345"
		break
	case "02":
		p.SaleBaseAmt = StandardValue
		p.DiscountBaseAmt = DiscountValue
		p.EventNo = "12345"
		break
	case "03":
		p.DiscountRate = 100 - DiscountValue
		p.EventNo = "12345"
		break
	case "B":
		p.SaleBaseAmt = StandardValue
		p.DiscountRate = 100 - DiscountValue
		break
	case "V":
		p.SaleBaseAmt = StandardValue
		p.DiscountRate = 100 - DiscountValue
		break
	case "P":
		p.SaleBaseAmt = StandardValue
		p.DiscountBaseAmt = DiscountValue
		break
	case "G", "R", "M":
		break
	default:
		break
	}
	return
}
