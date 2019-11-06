package promotion

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"time"
)

type PostCouponEvent struct {
	Id               int64     `json:"id"`
	BrandCode        string    `json:"brandCode"`
	EventTypeCode    string    `json:"eventTypeCode"`
	EventNo          string    `json:"eventNo" xorm:"index"`
	EventName        string    `json:"eventName"`
	EventDescription string    `json:"eventDescription"`
	StartDate        time.Time `json:"startDate"`
	InUserID         string    `json:"inUserId"`
	EndDate          time.Time `json:"endDate"`
	CreatedAt        time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt        time.Time `json:"updatedAt" xorm:"updated"`
}

func (PostCouponEvent) GetAll(ctx context.Context, v SearchInput) ([]PostCouponEvent, int64, error) {
	var (
		list       []PostCouponEvent
		totalCount int64
		err        error
	)
	query := factory.SaleRecordDB(ctx)
	if v.Q != "" {
		query = query.Where("event_no LIKE ?", v.Q+"%")
	}

	if v.BrandCode != "" {
		query = query.Where("brand_code = ?", v.BrandCode)
	}

	if err = setSortOrder(query, v.Sortby, v.Order, "post_coupon_event"); err != nil {
		return nil, 0, nil
	}
	totalCount, err = query.Limit(v.MaxResultCount, v.SkipCount).FindAndCount(&list)
	return list, totalCount, nil
}
