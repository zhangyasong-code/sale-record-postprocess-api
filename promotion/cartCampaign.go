package promotion

import (
	"context"
	"errors"
	offer "nomni/offer-api/models"
	"time"
)

type AggregatorType string

const (
	AggregatorTypeAll AggregatorType = "all"
	AggregatorTypeOne AggregatorType = "one"
)

type ComparerType string

const (
	ComparerTypeInclude          ComparerType = "in"
	ComparerTypeNotInclude       ComparerType = "nin"
	ComparerTypeGreaterThanEqual ComparerType = "gte"
	ComparerTypeLessThanEqual    ComparerType = "lte"
)

type StandardType string

const (
	StandardTypeNone  StandardType = "none"
	StandardTypePrice StandardType = "price"
	StandardTypeQty   StandardType = "qty"
)

type DiscountType string

const (
	DiscountTypeToFixedPrice  DiscountType = "to_fixed_price"  // Repeate as default, available in catalog
	DiscountTypeByFixedPrice  DiscountType = "by_fixed_price"  // Repeate as default, available in catalog
	DiscountTypeToFixedAmount DiscountType = "to_fixed_amount" // No repeate
	DiscountTypeToStepAmount  DiscountType = "to_step_amount"  // Repeate
	DiscountTypeByFixedAmount DiscountType = "by_fixed_amount" // No repeate
	DiscountTypeByStepAmount  DiscountType = "by_step_amount"  // Repeate
	DiscountTypeToPercentage  DiscountType = "to_percentage"   // -, available in catalog
	DiscountTypeByFixedQty    DiscountType = "by_fixed_qty"    // No repeate
	DiscountTypeByStepQty     DiscountType = "by_step_qty"     // Repeate
)

type AdditionalType string

const (
	AdditionalTypeGift    AdditionalType = "gift"
	AdditionalTypeMileage AdditionalType = "mileage"
)

type CartCampaign struct {
	Id             int64              `json:"id"`
	Code           string             `json:"code"`
	Name           string             `json:"name"`
	Desc           string             `json:"desc"`
	FeeRate        float64            `json:"feeRate"`
	Channels       []ChannelCondition `json:"channels"`
	StartAt        time.Time          `json:"startAt"`
	EndAt          time.Time          `json:"endAt"`
	FinalAt        time.Time          `json:"finalAt"` // 延期后的最终结束时间(== CSL：SaleEventEndDate + ExtendSalePermitDateCount)
	RulesetGroupId int64              `json:"rulesetGroupId"`
	Enable         bool               `json:"enable"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}

type ChannelCondition struct {
	CampaignId int64  `json:"-"`
	Type       string `json:"type"`
	Value      string `json:"value"`
}

type CartRulesetGroup struct {
	Id             int64               `json:"id"`
	TenantCode     string              `json:"-"`
	TemplateCode   string              `json:"templateCode"`
	Type           OfferType           `json:"type"`
	AdditionalType AdditionalType      `json:"additionalType"`
	Name           string              `json:"name"`
	Rulesets       []CartRuleset       `json:"rulesets"`
	Customers      []CustomerCondition `json:"customers"`
	Actions        []CartAction        `json:"actions"`
	Enable         bool                `json:"enable"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
}

type CartRuleset struct {
	Id         int64           `json:"id"`
	GroupId    int64           `json:"-" `
	Enable     bool            `json:"enable"`
	Qty        int64           `json:"qty"`
	Conditions []CartCondition `json:"conditions,omitempty"`
}

type CartCondition struct {
	RulesetId int64        `json:"-"`
	Type      string       `json:"type"  swagger:"enum(brand_code|product_id|list_price)"`
	Comparer  ComparerType `json:"comparer" swagger:"enum(in|nin|gte|lte)"`
	Targets   []CartTarget `json:"targets"`
}

type CartTarget struct {
	RulesetId     int64  `json:"-" `
	ConditionType string `json:"-" `
	Value         string `json:"value" `
}

type CustomerCondition struct {
	RulesetGroupId int64            `json:"-"`
	SeqNo          int              `json:"seqNo"`
	Type           string           `json:"type" swagger:"enum(member_id|grade_id|birthday_month)"`
	Comparer       ComparerType     `json:"comparer" ` // in、nin
	Cnt            int              `json:"-" `
	Targets        []CustomerTarget `json:"targets"`
}

type CustomerTarget struct {
	RulesetGroupId int64  `json:"-"`
	ConditionSeqNo int    `json:"-"`
	ConditionType  string `json:"-"`
	Value          string `json:"value"`
}

type CartAction struct {
	StandardType    StandardType         `json:"standardType" swagger:"enum(none|price|qty)"`
	StandardValue   float64              `json:"standardValue"`
	DiscountType    DiscountType         `json:"discountType" swagger:"enum(by_fixed_amount|to_fixed_amount|by_fixed_price|to_fixed_price|to_percetage|by_fixed_qty)"`
	DiscountValue   float64              `json:"discountValue"`
	DiscountTargets []CartDiscountTarget `json:"discountTargets,omitempty"`
}

type CartDiscountTarget struct {
	Type  string `json:"type" swagger:"enum(product_id|sku_id)"`
	Value int64  `json:"value"`
}

func CartToCSLEvent(ctx context.Context, c CartCampaign, ruleGroup *CartRulesetGroup) (*PromotionEvent, error) {
	var brandCode, storeCode string
	eventType, err := ToCSLOfferType(ruleGroup.Type, ruleGroup.TemplateCode)
	if err != nil {
		return nil, err
	}
	brand, store, err := getBrandAndStore(ctx, c.Channels)
	if err != nil {
		return nil, err
	}

	if brand == nil {
		return nil, errors.New("brand is not exist")
	}
	brandCode = brand.Code
	if store == nil {
		storeCode = ""
	} else {
		storeCode = store.Code
	}

	promotion := &Promotion{
		EventType: eventType,
	}
	if len(ruleGroup.Actions) > 0 {
		promotion.ToCSLDisCount(ruleGroup.Actions[0].StandardValue, ruleGroup.Actions[0].DiscountValue)
	}
	offerNo := offer.NewOfferNo(offer.CampaignTypeCart, c.Id, ruleGroup.Id)
	return &PromotionEvent{
		OfferNo:                   string(offerNo),
		BrandCode:                 brandCode,
		ShopCode:                  storeCode,
		EventTypeCode:             eventType,
		EventName:                 c.Name,
		EventDescription:          c.Desc,
		StartDate:                 c.StartAt,
		EndDate:                   c.FinalAt,
		ExtendSalePermitDateCount: 0,
		NormalSaleRecognitionChk:  promotion.NormalSaleRecognitionChk,
		FeeRate:                   c.FeeRate,
		InUserID:                  "mslv2.0",
		SaleBaseAmt:               promotion.SaleBaseAmt,
		DiscountBaseAmt:           promotion.DiscountBaseAmt,
		DiscountRate:              promotion.DiscountRate,
		StaffSaleChk:              false,
	}, nil
}
