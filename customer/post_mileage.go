package customer

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"nhub/sale-record-postprocess-api/models"
	"time"

	"github.com/go-xorm/xorm"
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
	Id                  int64      `json:"id"`
	TenantCode          string     `json:"tenantCode" xorm:"index VARCHAR(50) notnull" validate:"required"`
	StoreId             int64      `json:"storeId" xorm:"index notnull" validate:"gte=0"`
	CustomerId          int64      `json:"customerId" xorm:"index notnull" validate:"gte=0"`
	BrandId             int64      `json:"brandId" xorm:"notnull"`
	GradeId             int64      `json:"gradeId"`
	TransactionId       int64      `json:"transactionId" xorm:"index default 0" validate:"required"`
	OrderId             int64      `json:"orderId" xorm:"index default 0" validate:"required"`
	RefundId            int64      `json:"refundId" xorm:"index default 0"`
	UseType             UseType    `json:"useType" xorm:"VARCHAR(25)"`
	Point               float64    `json:"point" xorm:"decimal(19,2)"`
	PointPrice          float64    `json:"pointPrice" xorm:"decimal(19,2)"`
	CustMileagePolicyNo int64      `json:"custMileagePolicyNo"`
	CreateAt            *time.Time `json:"createAt" xorm:"created"`
	UpdateAt            *time.Time `json:"updateAt" xorm:"updated"`
}

type PostMileageDtl struct {
	Id                  int64      `json:"id"`
	PostMileageId       int64      `json:"postMileageId"`
	TransactionDtlId    int64      `json:"transactionDtlId" xorm:"index default 0"`
	OrderItemId         int64      `json:"orderItemId" xorm:"index default 0"`
	RefundItemId        int64      `json:"refundItemId" xorm:"index default 0"`
	UseType             UseType    `json:"useType" xorm:"VARCHAR(25)"`
	Point               float64    `json:"point" xorm:"decimal(19,2)"`
	PointPrice          float64    `json:"pointPrice" xorm:"decimal(19,2)"`
	CustMileagePolicyNo int64      `json:"custMileagePolicyNo"`
	CreateAt            *time.Time `json:"createAt" xorm:"created"`
	UpdateAt            *time.Time `json:"updateAt" xorm:"updated"`
}

func (PostMileage) MakePostMileage(ctx context.Context, mileage Mileage, record models.SaleRecordEvent) (*PostMileage, error) {
	mileageMall, err := MileageMall{}.GetMembershipGrade(ctx, mileage.MallId, mileage.MemberId, mileage.TenantCode)
	if err != nil {
		return nil, err
	}
	custMileagePolicy, err := CustMileagePolicy{}.GetCustMileagePolicy(ctx, record.AssortedSaleRecordDtlList[0].BrandCode)
	if err != nil {
		return nil, err
	}

	var useType UseType
	if mileage.Type == CONSUME {
		if record.RefundId != 0 {
			useType = UseTypeUsedCancel
		} else {
			useType = UseTypeUsed
		}
	} else {
		if record.RefundId != 0 {
			useType = UseTypeEarnCancel
		} else {
			useType = UseTypeEarn
		}
	}

	postMileage := &PostMileage{
		TenantCode:          mileage.TenantCode,
		StoreId:             mileage.StoreId,
		CustomerId:          mileage.MemberId,
		BrandId:             mileage.MallId,
		GradeId:             mileageMall.GradeId,
		TransactionId:       record.TransactionId,
		OrderId:             record.OrderId,
		RefundId:            record.RefundId,
		UseType:             useType,
		Point:               mileage.Point,
		PointPrice:          mileage.PointPrice,
		CustMileagePolicyNo: custMileagePolicy.CustMileagePolicyNo,
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

func (PostMileageDtl) MakePostMileageDtls(postMileage *PostMileage, mileageDtls []MileageDtl, recordDtls []models.AssortedSaleRecordDtl) []PostMileageDtl {
	var postMileageDtls []PostMileageDtl
	for _, mileageDtl := range mileageDtls {
		postMileageDtl := PostMileageDtl{
			PostMileageId:       postMileage.Id,
			UseType:             postMileage.UseType,
			Point:               mileageDtl.Point,
			PointPrice:          mileageDtl.PointPrice,
			CustMileagePolicyNo: postMileage.CustMileagePolicyNo,
		}
		for _, recordDtl := range recordDtls {
			if recordDtl.OrderItemId == mileageDtl.ItemId && recordDtl.RefundItemId == mileageDtl.PreItemId {
				postMileageDtl.TransactionDtlId = recordDtl.Id
				postMileageDtl.OrderItemId = recordDtl.OrderItemId
				postMileageDtl.RefundItemId = recordDtl.RefundItemId
				postMileageDtls = append(postMileageDtls, postMileageDtl)
			}
		}
	}
	return postMileageDtls
}
func (PostMileageDtl) CreateBatch(ctx context.Context, v []PostMileageDtl) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(v); err != nil {
		return err
	}
	return nil
}
func (PostMileageDtl) GetByKey(ctx context.Context, transactionDtlId, orderItemId, refundItemId int64, useType UseType) (bool, PostMileageDtl, error) {
	return PostMileageDtl{}.GetPostMileageDtl(ctx, transactionDtlId, orderItemId, refundItemId, useType)
}

func (PostMileageDtl) GetPostMileageDtl(ctx context.Context, transactionDtlId, orderItemId, refundItemId int64, useType UseType) (bool, PostMileageDtl, error) {
	res := PostMileageDtl{}
	query := func() xorm.Interface {
		q := factory.SaleRecordDB(ctx).Where("1 = 1")
		if orderItemId != 0 {
			q.And("order_item_id = ?", orderItemId)
		}
		if useType != UseTypeUnKnown {
			q.And("use_type = ?", string(useType))
		}
		if transactionDtlId != 0 {
			q.And("transaction_dtl_id = ?", transactionDtlId)
		}
		if refundItemId != 0 {
			q.And("refund_item_id = ?", refundItemId)
		}
		return q
	}
	exist, err := query().Get(&res)
	if err != nil {
		return exist, PostMileageDtl{}, err
	}
	if !exist {
		return exist, PostMileageDtl{}, nil
	}
	return exist, res, nil
}
