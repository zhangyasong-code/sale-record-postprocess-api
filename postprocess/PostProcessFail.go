package postprocess

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"time"

	"github.com/go-xorm/xorm"
)

type ModuleType string

const (
	ModuleUnknown    ModuleType = ""
	ModulePromotion  ModuleType = "Promotion"
	ModuleMileage    ModuleType = "Mileage"
	ModulePay        ModuleType = "Pay"
	ModuleSalePerson ModuleType = "SalePerson"
	ModuleSaleFee    ModuleType = "SaleFee"
)

type PostProcessSuccess struct {
	Id           int64     `json:"id"`
	ModuleType   string    `json:"moduleType" xorm:"index VARCHAR(50)"`
	OrderId      int64     `json:"orderId" xorm:"index default 0" validate:"required"`
	RefundId     int64     `json:"refund" xorm:"index default 0"`
	IsSuccess    bool      `json:"isSuccess" xorm:"index notnull default false"`
	Error        string    `json:"error" xorm:"VARCHAR(1000)" validate:"required"`
	ModuleEntity string    `json:"moduleEntity" xorm:"VARCHAR(5000)"`
	CreatedAt    time.Time `json:"createdAt" xorm:"index created"`
	UpdatedAt    time.Time `json:"updatedAt" xorm:"updated"`
}

type PostFailParam struct {
	OrderId        int64      `json:"orderId"`
	RefundId       int64      `json:"refundId"`
	ModuleType     ModuleType `json:"moduleType"`
	IsSuccess      bool       `json:"isSuccess"`
	MaxResultCount int        `json:"maxResultCount"`
	SkipCount      int        `json:"skipCount"`
}

func (post *PostProcessSuccess) Save(ctx context.Context) error {
	var postProcessSuccess PostProcessSuccess
	has, err := factory.SaleRecordDB(ctx).Where("order_id = ?", post.OrderId).And("refund_id = ?", post.RefundId).And("module_type = ?", post.ModuleType).Get(&postProcessSuccess)
	if err != nil {
		return err
	}
	if !has {
		if _, err := factory.SaleRecordDB(ctx).Insert(post); err != nil {
			return err
		}
	} else {
		if err := post.Update(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (post *PostProcessSuccess) Update(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Where("order_id = ?", post.OrderId).And("refund_id = ?", post.RefundId).And("module_type = ?", post.ModuleType).AllCols().Update(post); err != nil {
		return err
	}
	return nil
}

func (PostProcessSuccess) GetAll(ctx context.Context, postFailParam PostFailParam) (int64, []PostProcessSuccess, error) {
	var postProcessSuccess []PostProcessSuccess
	query := func() xorm.Interface {
		query := factory.SaleRecordDB(ctx).Where("1 = 1").And("isSuccess = ?", postFailParam.IsSuccess)
		if postFailParam.OrderId != 0 {
			query.And("order_id = ?", postFailParam.OrderId)
		}
		if postFailParam.RefundId != 0 {
			query.And("refund_id = ?", postFailParam.RefundId)
		}
		if postFailParam.ModuleType != ModuleUnknown {
			query.And("module_type = ?", string(postFailParam.ModuleType))
		}
		return query
	}
	totalCount, err := query().Limit(postFailParam.MaxResultCount, postFailParam.SkipCount).FindAndCount(&postProcessSuccess)
	if err != nil {
		return 0, nil, err
	}
	return totalCount, postProcessSuccess, nil
}

func IsSuccessd(ctx context.Context) error {
	_, postProcessSuccess, err := PostProcessSuccess{}.GetAll(ctx, PostFailParam{IsSuccess: false})
	if err != nil {
		return err
	}
	for _, postProcess := range postProcessSuccess {
		postProcess.IsSuccess = true
		if err := postProcess.Update(ctx); err != nil {
			return err
		}
	}

	return nil
}
