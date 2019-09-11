package promotion

import (
	"context"
	"errors"
	"strconv"
)

type ConditionType string

const (
	ConditionTypeStoreId = "store_id"
	ConditionTypeBrandId = "brand_id"
)

type Brand struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Store struct {
	Id     int64   `json:"id"`
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Brands []Brand `json:"brands"`
	Remark Remark  `json:"remark"`
}

type Remark struct {
	ElandShopInfos []ElandShopInfo `json:"elandShopInfos"`
}

type ElandShopInfo struct {
	BrandCode string `json:"brandCode"`
	BrandId   int64  `json:"brandId"`
	IsCheif   bool   `json:"isChief"`
	ShopCode  string `json:"shopCode"`
}

type Promotion struct {
	EventType                string  `json:"eventType"`
	SaleBaseAmt              float64 `json:"saleBaseAmt"`
	DiscountBaseAmt          float64 `json:"discountBaseAmt"`
	DiscountRate             float64 `json:"discountRate"`
	NormalSaleRecognitionChk bool    `json:"normalSaleRecognitionChk"`
}

func getBrandAndStore(ctx context.Context, channels []ChannelCondition) (*Brand, []Store, error) {
	var (
		brandId  int64
		storeIds []string
		err      error
		brand    *Brand
	)
	for i := range channels {
		if channels[i].Type == ConditionTypeBrandId {
			brandId, err = strconv.ParseInt(channels[i].Value, 10, 64)
			if err != nil {
				return nil, nil, errors.New("brandId invalid")
			}
		}
		if channels[i].Type == ConditionTypeStoreId {
			storeIds = append(storeIds, channels[i].Value)
		}
	}
	if brandId != 0 {
		brand, err = getBrand(ctx, brandId)
		if err != nil {
			return nil, nil, err
		}
		return brand, nil, nil
	}

	if len(storeIds) == 0 {
		return nil, nil, errors.New("brand and store not exist")
	}
	stores, err := getStores(ctx, storeIds)
	if err != nil {
		return nil, nil, err
	}
	return brand, stores, nil
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
	//TODO:区分优惠券和custEvent模板
	case "B1", "B2", "C1", "C2", "C6", "D1":
		if offerType == OfferTypeBrand {
			eventTypeCode = "02"
		}
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
		p.NormalSaleRecognitionChk = true
		break
	case "02":
		p.SaleBaseAmt = StandardValue
		p.DiscountBaseAmt = DiscountValue
		p.NormalSaleRecognitionChk = false
		break
	case "03":
		p.DiscountRate = 100 - DiscountValue
		p.NormalSaleRecognitionChk = false
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
