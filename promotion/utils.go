package promotion

import (
	"context"
	"errors"
	"strconv"
)

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
