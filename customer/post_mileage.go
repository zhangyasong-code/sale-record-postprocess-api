package customer

import (
	"context"
	"fmt"
	"nhub/sale-record-postprocess-api/factory"
	"nhub/sale-record-postprocess-api/models"
	"time"
)

type UseType string

const (
	UseTypeUnKnown    UseType = ""
	UseTypeEarn       UseType = "Earn"
	UseTypeEarnCancel UseType = "EarnCancel"
	UseTypeUsed       UseType = "Used"
	UseTypeUsedCancel UseType = "UsedCancel"
)

type PostMileage struct {
	Id            int64      `json:"id"`
	TenantCode    string     `json:"tenantCode" xorm:"index VARCHAR(50) notnull" validate:"required"`
	StoreId       int64      `json:"storeId" xorm:"index notnull" validate:"gte=0"`
	CustomerId    int64      `json:"customerId" xorm:"index notnull" validate:"gte=0"`
	BrandId       int64      `json:"brandId" xorm:"notnull"`
	BrandCode     string     `json:"brandCode" xorm:"notnull"`
	GradeId       int64      `json:"gradeId"`
	TransactionId int64      `json:"transactionId" xorm:"index default 0" validate:"required"`
	OrderId       int64      `json:"orderId" xorm:"index default 0" validate:"required"`
	RefundId      int64      `json:"refundId" xorm:"index default 0"`
	CreateAt      *time.Time `json:"createAt" xorm:"created"`
	UpdateAt      *time.Time `json:"updateAt" xorm:"updated"`
}

func (PostMileage) MakePostMileage(ctx context.Context, record models.SaleRecordEvent) (*PostMileage, error) {
	/*获取品牌Id*/
	brandId, err := Store{}.GetBrandIdByStoreId(ctx, record.StoreId)
	if err != nil {
		return nil, err
	}
	if brandId == 0 {
		return nil, fmt.Errorf("Fail to get brandId")
	}

	mileageMall, err := MileageMall{}.GetMembershipGrade(ctx, brandId, record.CustomerId, record.TenantCode)
	if err != nil {
		return nil, err
	}

	brandCode, err := BrandInfo{}.GetBrandInfo(ctx, brandId)
	if err != nil {
		return nil, err
	}

	postMileage := &PostMileage{
		TenantCode:    record.TenantCode,
		StoreId:       record.StoreId,
		CustomerId:    record.CustomerId,
		BrandId:       brandId,
		BrandCode:     brandCode,
		GradeId:       mileageMall.GradeId,
		TransactionId: record.TransactionId,
		OrderId:       record.OrderId,
		RefundId:      record.RefundId,
	}
	return postMileage, nil
}

func (o *PostMileage) Create(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(o); err != nil {
		return err
	}
	return nil
}

func (PostMileage) CheckOrderRefundExist(ctx context.Context, transactionId int64) (bool, error) {
	postMileage := PostMileage{}
	has, err := factory.SaleRecordDB(ctx).Table("post_mileage").
		Where("transaction_id=?", transactionId).Get(&postMileage)
	if err != nil {
		return true, err
	}
	return has, nil
}
