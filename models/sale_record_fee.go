package models

import (
	"context"
	"errors"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"nhub/sale-record-postprocess-api/factory"
	"nhub/sale-record-postprocess-api/promotion"
	"strconv"
	"strings"
	"time"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

type PostSaleRecordFee struct {
	Id                     int64     `json:"id" query:"id"`
	TransactionId          string    `json:"transactionId" query:"transactionId" xorm:"index VARCHAR(30) notnull" validate:"required"`
	SaleRecordDtlId        int64     `json:"saleRecordDtlId" query:"saleRecordDtlId" xorm:"index VARCHAR(30) notnull" validate:"required"`
	SaleRecordOfferId      int64     `json:"saleRecordOfferId" query:"saleRecordOfferId" xorm:"index notnull" validate:"gte=0"`
	OrderId                int64     `json:"orderId" query:"orderId" xorm:"index notnull" validate:"gte=0"`
	OrderItemId            int64     `json:"orderItemId" query:"orderItemId" xorm:"index notnull" validate:"gte=0"`
	RefundId               int64     `json:"refundId" query:"refundId" xorm:"index notnull" validate:"gte=0"`
	RefundItemId           int64     `json:"refundItemId" query:"refundItemId" xorm:"index notnull" validate:"gte=0"`
	CustomerId             int64     `json:"customerId" query:"customerId" xorm:"index notnull" validate:"gte=0"`
	StoreId                int64     `json:"storeId" query:"storeId" xorm:"index notnull" validate:"gte=0"`
	EventType              string    `json:"eventType" query:"eventType" xorm:"VARCHAR(20)"`
	TotalSalePrice         float64   `json:"totalSalePrice" query:"totalSalePrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	TotalPaymentPrice      float64   `json:"totalPaymentPrice" query:"totalPaymentPrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	Mileage                float64   `json:"mileage" query:"mileage" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	MileagePrice           float64   `json:"mileagePrice" query:"mileagePrice" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	ContractFeeRate        float64   `json:"contractFeeRate" query:"contractFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	EventFeeRate           float64   `json:"eventFeeRate" query:"eventFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	AppliedFeeRate         float64   `json:"appliedFeeRate" query:"appliedFeeRate" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	FeeAmount              float64   `json:"feeAmount" query:"feeAmount" xorm:"DECIMAL(18,2) default 0" validate:"gte=0"`
	TransactionChannelType string    `json:"transactionChannelType" query:"transactionChannelType" xorm:"index VARCHAR(20) notnull"`
	CreatedAt              time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt              time.Time `json:"updatedAt" xorm:"updated"`
}

type PostFailCreateSaleFee struct {
	Id            int64     `json:"id" query:"id"`
	TransactionId string    `json:"transactionId" query:"transactionId" xorm:"index VARCHAR(30) notnull" validate:"required"`
	IsProcessed   bool      `json:"isProcessed" query:"isProcessed" xorm:"index notnull default false" validate:"required"`
	CreatedAt     time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt     time.Time `json:"updatedAt" xorm:"updated"`
}

type ResultShopMessage struct {
	Success bool `json:"success"`
	Result  struct {
		Items []struct {
			Contracts []Contract `json:"contracts"`
		} `json:"items"`
	} `json:"result"`
	Error struct{} `json:"error"`
}

type Contract struct {
	Id         int64     `json:"id"`
	StoreId    int64     `json:"storeId"`
	BrandId    int64     `json:"brandId"`
	ContractNo string    `json:"contractNo"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	BaseRate   float64   `json:"baseRate"`
}

func (p *PostSaleRecordFee) Save(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(p); err != nil {
		return err
	}
	return nil
}

func (pf *PostFailCreateSaleFee) Save(ctx context.Context) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(pf); err != nil {
		return err
	}
	return nil
}

func getContracts(ctx context.Context, storeId int64) (*ResultShopMessage, error) {
	var result ResultShopMessage
	if _, err := httpreq.New(http.MethodGet, config.Config().
		Services.PlaceManagementApi+"/v1/store/getallinfo?skipCount=0&maxResultCount=10&withContract=true&propsEnable=true&storeIds="+strconv.FormatInt(storeId, 10), nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&result); err != nil {
		return nil, err
	}
	if !result.Success {
		logrus.WithFields(logrus.Fields{"ServiceName": "PlaceManagementApi", "errorMessage": result.Error, "storeId": storeId}).Error("GetContracts failed!")
		return nil, errors.New("GetContracts failed!")
	}
	return &result, nil
}

func (PostSaleRecordFee) GetContractFeeRate(ctx context.Context, storeId, brandId int64, transactionCreateDate time.Time) (float64, error) {
	var contractFeeRate float64
	// Use the storeId to query contracts
	resultShopMessage, err := getContracts(ctx, storeId)
	if err != nil {
		return 0, err
	}
	if len(resultShopMessage.Result.Items) != 0 {
		for _, contract := range resultShopMessage.Result.Items[0].Contracts {
			if contract.StartDate.Before(transactionCreateDate) && contract.EndDate.After(transactionCreateDate) && contract.BrandId == brandId {
				contractFeeRate = contract.BaseRate
				break
			}
		}
	}
	return contractFeeRate, nil
}

func (PostSaleRecordFee) GetPromotionEvent(ctx context.Context, offerNo string) (*promotion.PromotionEvent, error) {
	arr := strings.Split(offerNo, "-")
	campaignId, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return nil, err
	}
	ruleId, err := strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		return nil, err
	}
	promotionEvent, err := promotion.GetById(ctx, campaignId, ruleId)
	if err != nil {
		return nil, err
	}
	return promotionEvent, nil
}

func (PostSaleRecordFee) GetPostMileageDtl(ctx context.Context, orderItemId, refundItemId int64) (*PostMileageDtl, error) {
	var o PostMileageDtl
	exist, err := factory.SaleRecordDB(ctx).Where("order_item_id = ?", orderItemId).And("refund_item_id = ?", refundItemId).Get(&o)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, errors.New("OrgMileageDtl is not exist")
	}
	return &o, nil
}
