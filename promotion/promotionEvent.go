package promotion

import (
	"context"
	"errors"
	"nhub/sale-record-postprocess-api/factory"
	"time"
)

type PromotionEvent struct {
	OfferNo                   string    `json:"offerNo" xorm:"pk"`
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
