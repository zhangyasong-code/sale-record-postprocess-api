package promotion

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	offer "nomni/offer-api/models"

	"github.com/go-xorm/xorm"
)

type ConditionType string

const (
	ConditionTypeStoreId = "store_id"
	ConditionTypeBrandId = "brand_id"
)

type SearchInput struct {
	Q              string   `query:"q"`
	BrandCode      string   `query:"brandCode"`
	EventTypeCodes string   `query:"eventTypeCodes"`
	Sortby         []string `query:"sortby"`
	Order          []string `query:"order"`
	SkipCount      int      `query:"skipCount"`
	MaxResultCount int      `query:"maxResultCount"`
}

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
	case "A2", "A2-O":
		if offerType == OfferTypeBrand {
			eventTypeCode = "03"
		} else if offerType == OfferTypeMember {
			eventTypeCode = "B"
		} else if offerType == OfferTypeChannel {
			eventTypeCode = "V"
		}
	case "B1", "B2", "C1", "C6", "D1":
		if offerType == OfferTypeBrand {
			eventTypeCode = "02"
		}
	case "C5", "C2":
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
		p.SaleBaseAmt, p.DiscountBaseAmt = getDefaultValue(StandardValue, DiscountValue)
		p.NormalSaleRecognitionChk = true
		break
	case "02":
		p.SaleBaseAmt, p.DiscountBaseAmt = getDefaultValue(StandardValue, DiscountValue)
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

func GetFeeRate(offerType OfferType, simulations []CampaignSimulation) float64 {
	if offerType == OfferTypeBrand && len(simulations) > 0 {
		return simulations[0].SaleEventFeeRate
	} else {
		return 0
	}
}

//offer金额小于10，或是类型为数量时，默认传10
func getDefaultValue(standardValue, discountValue float64) (float64, float64) {
	if discountValue < 10 {
		standardValue = 11
		discountValue = 10
	} else if standardValue <= discountValue {
		standardValue = discountValue + 1
	}
	return standardValue, discountValue
}

func setSortOrder(q xorm.Interface, sortby, order []string, table ...string) error {
	connect := func(col string) string {
		if len(table) > 0 {
			return table[0] + "." + col
		}
		return col
	}

	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				v = connect(v)
				if order[i] == "desc" {
					q.Desc(v)
				} else if order[i] == "asc" {
					q.Asc(v)
				} else {
					return errors.New("Invalid order. Must be either [asc|desc]")
				}
			}
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				v = connect(v)
				if order[0] == "desc" {
					q.Desc(v)
				} else if order[0] == "asc" {
					q.Asc(v)
				} else {
					return errors.New("Invalid order. Must be either [asc|desc]")
				}
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return errors.New("'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return errors.New("unused 'order' fields")
		}
	}
	return nil
}

func convertOfferNo(no offer.OfferNo) []offer.OfferNo {
	if no.Version() == 0 {
		return []offer.OfferNo{no}
	} else {
		return []offer.OfferNo{
			offer.OfferNo(fmt.Sprintf("%d-%d-%d", no.CampaignType(), no.CampaignId(), no.RulesetOrGroupId())),
			offer.NewOfferNo(no.CampaignType(), no.CampaignId(), no.RulesetOrGroupId()),
		}
	}
}
