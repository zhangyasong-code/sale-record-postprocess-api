package promotion

import (
	"context"
	"errors"
	"nhub/sale-record-postprocess-api/factory"
	offer "nomni/offer-api/models"
	"strings"
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
	ErrorMsg                  string    `json:"errorMsg"`
	CreatedAt                 time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt                 time.Time `json:"updatedAt" xorm:"updated"`
}

func (p *PromotionEvent) create(ctx context.Context) error {
	_, err := factory.SaleRecordDB(ctx).Insert(p)
	return err
}

func GetByNo(ctx context.Context, no string) (*PromotionEvent, error) {
	//将新结构的offerNo转换成旧结构
	nos := convertOfferNo(offer.OfferNo(no))
	var p PromotionEvent
	exist, err := factory.SaleRecordDB(ctx).In("offer_no", nos).Desc("created_at").Get(&p)
	if err != nil {
		return nil, err
	} else if !exist {
		uploadErr := reUploadOffer(ctx, no)
		if uploadErr != nil {
			return nil, errors.New("promotionEvent offer_no = '" + no + "' is not exist")
		}
		return nil, errors.New("promotionEvent offer_no = '" + no + "' is not exist")
	}
	if p.ErrorMsg != "" {
		_ = reUploadOffer(ctx, no)
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
	for i := range stores {
		p.ShopCode = stores[i].Code
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

func (p PromotionEvent) createOrUpdate(ctx context.Context) error {
	if p.EventTypeCode == "" {
		return errors.New("PromotionEvent eventTypeCode is null")
	}
	if p.FeeRate <= 0 && (p.EventTypeCode == "01" || p.EventTypeCode == "02" || p.EventTypeCode == "03") {
		return errors.New("PromotionEvent feeRate lower then equal 0")
	}
	var promotion PromotionEvent
	exist, err := factory.SaleRecordDB(ctx).Where("offer_no = ?", p.OfferNo).Get(&promotion)
	if err != nil {
		return err
	}
	if exist {
		_, err := factory.SaleRecordDB(ctx).ID(promotion.Id).Cols("event_no,event_type_code,sale_base_amt,discount_base_amount,discount_rate,fee_rate,error_msg").Update(p)
		if err != nil {
			return err
		}
	} else {
		if err := p.create(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (PromotionEvent) GetAll(ctx context.Context, v SearchInput) ([]PromotionEvent, int64, error) {
	var (
		list       []PromotionEvent
		totalCount int64
		err        error
	)
	query := factory.SaleRecordDB(ctx)
	if v.Q != "" {
		query = query.Where("event_name LIKE ? OR offer_no LIKE ? OR event_no LIKE ?", v.Q+"%", v.Q+"%", v.Q+"%")
	}
	if v.EventTypeCodes != "" {
		codes := strings.Split(v.EventTypeCodes, ",")
		query = query.In("event_type_code", codes)
	}
	if err = setSortOrder(query, v.Sortby, v.Order, "promotion_event"); err != nil {
		return nil, 0, nil
	}
	totalCount, err = query.Limit(v.MaxResultCount, v.SkipCount).FindAndCount(&list)
	return list, totalCount, nil
}
