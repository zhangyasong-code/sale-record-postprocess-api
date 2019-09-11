package promotion

import (
	"context"
	"errors"
	"nhub/sale-record-postprocess-api/factory"
	"strconv"
	"time"
)

type PromotionEvent struct {
	Id                        int64     `json:"id" xorm:"pk"`
	OfferNo                   string    `json:"offerNo" xorm:"index"`
	BrandCode                 string    `json:"brandCode"`
	ShopCode                  string    `json:"shopCode"` //是否需要（SaleEvent）
	EventTypeCode             string    `json:"eventTypeCode"`
	EventName                 string    `json:"eventName"`
	EventNo                   string    `json:"eventNo"`
	EventDescription          string    `json:"eventDescription"`
	StartDate                 time.Time `json:"startDate"`
	EndDate                   time.Time `json:"endDate"`
	ExtendSalePermitDateCount int       `json:"extendSalePermitDateCount"` //扩展天数
	NormalSaleRecognitionChk  bool      `json:"normalSaleRecognitionChk"`  //活动销售额是否正常
	FeeRate                   float64   `json:"feeRate"`
	InUserID                  string    `json:"inUserId"`
	SaleBaseAmt               float64   `json:"saleBaseAmt"`
	DiscountBaseAmt           float64   `json:"discountBaseAmt"`
	DiscountRate              float64   `json:"discountRate"`
	StaffSaleChk              bool      `json:"staffSaleChk"`
	CreatedAt                 time.Time `json:"createdAt" xorm:"created"`
}

func (p *PromotionEvent) create(ctx context.Context) error {
	_, err := factory.SaleRecordDB(ctx).Insert(p)
	return err
}

func GetByNo(ctx context.Context, no string) (*PromotionEvent, error) {
	var p PromotionEvent
	exist, err := factory.SaleRecordDB(ctx).Where("offer_no = ?", no).Get(&p)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, errors.New("promotionEvent is not exist")
	}
	return &p, nil
}

func (p PromotionEvent) CreateByStoreOrBrand(brand *Brand, stores []Store, channels []ChannelCondition) []PromotionEvent {
	var list []PromotionEvent
	if brand != nil {
		p.BrandCode = brand.Code
		list = append(list, p)
		return list
	}
	getFeeRate := func(id int64) float64 {
		for _, channel := range channels {
			if strconv.FormatInt(id, 10) == channel.Value && channel.Type == ConditionTypeStoreId {
				return channel.FeeRate
			}
		}
		return 0
	}
	for i := range stores {
		p.ShopCode = stores[i].Code
		p.FeeRate = getFeeRate(stores[i].Id)
		for _, info := range stores[i].Remark.ElandShopInfos {
			if info.IsCheif {
				p.BrandCode = info.BrandCode

			}
		}
		list = append(list, p)
	}
	return list
}

func (PromotionEvent) createInArrary(ctx context.Context, promotions []PromotionEvent) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(&promotions); err != nil {
		return err
	}
	return nil
}
