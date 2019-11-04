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
	Id            int64     `json:"id"`
	ModuleType    string    `json:"moduleType" xorm:"index VARCHAR(50)"`
	TransactionId int64     `json:"transactionId" xorm:"index default 0" validate:"required"`
	OrderId       int64     `json:"orderId" xorm:"index default 0" validate:"required"`
	RefundId      int64     `json:"refund" xorm:"index default 0"`
	IsSuccess     bool      `json:"isSuccess" xorm:"index notnull default false"`
	Error         string    `json:"error" xorm:"VARCHAR(1000)" validate:"required"`
	ModuleEntity  string    `json:"moduleEntity" xorm:"TEXT"`
	CreatedAt     time.Time `json:"createdAt" xorm:"index created"`
	UpdatedAt     time.Time `json:"updatedAt" xorm:"updated"`
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

func (PostProcessSuccess) Get(ctx context.Context, isSuccess bool, transactionId int64, moduleType string) (PostProcessSuccess, error) {
	var postProcessSuccess PostProcessSuccess
	exist, err := factory.SaleRecordDB(ctx).Where("1 = 1").And("is_Success = ?", isSuccess).And("transaction_id = ?", transactionId).And("module_type = ?", moduleType).Get(&postProcessSuccess)
	if err != nil {
		return postProcessSuccess, err
	} else if !exist {
		return postProcessSuccess, nil
	}

	return postProcessSuccess, nil
}

func (PostProcessSuccess) GetAll(ctx context.Context, isSuccess bool, transactionId, orderId, refundId int64, moduleType string, skipCount int, maxResultCount int) (int64, []PostProcessSuccess, error) {
	var postProcessSuccess []PostProcessSuccess
	query := func() xorm.Interface {
		query := factory.SaleRecordDB(ctx).Where("1 = 1").And("is_Success = ?", isSuccess)
		if transactionId != 0 {
			query.And("transaction_id = ?", transactionId)
		}
		if orderId != 0 {
			query.And("order_id = ?", orderId)
		}
		if refundId != 0 {
			query.And("refund_id = ?", refundId)
		}
		if moduleType != "" {
			query.And("module_type = ?", moduleType)
		}
		return query
	}
	if maxResultCount != 0 {
		query().Limit(maxResultCount, skipCount)
	}
	totalCount, err := query().Desc("id").FindAndCount(&postProcessSuccess)
	if err != nil {
		return 0, nil, err
	}
	return totalCount, postProcessSuccess, nil
}
