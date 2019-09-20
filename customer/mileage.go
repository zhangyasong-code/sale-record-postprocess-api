package customer

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"nomni/utils/auth"
	"time"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

const (
	SALE    string = "E"
	CONSUME string = "M"
	ADJUST  string = "A"
)

type Mileage struct {
	Id          int64        `json:"id,omitempty"`         /*ID 主键*/
	TenantCode  string       `json:"tenantCode,omitempty"` /*租户代码*/
	MallId      int64        `json:"mallId,omitempty"`     /*购物中心代码*/
	StoreId     int64        `json:"storeId,omitempty"`    /*店铺代码*/
	MemberId    int64        `json:"memberId"`             /*会员id*/
	Type        string       `json:"type,omitempty"`       /*类型*/
	ChannelId   int64        `json:"channelId"`            /*渠道*/
	TradeDate   time.Time    `json:"tradeDate,omitempty"`  /*交易日期*/
	TradeNo     string       `json:"tradeNo,omitempty"`    /*交易单号*/
	PreTradeNo  string       `json:"preTradeNo,omitempty"` /*原交易单号*/
	SaleAmount  float64      `json:"saleAmount"`           /*整单金额*/
	PayAmount   float64      `json:"payAmount"`            /*实付金额*/
	Point       float64      `json:"point"`                /*积分数量*/
	PointPrice  float64      `json:"pointPrice"`           /*积分抵扣金额*/
	Remark      string       `json:"remark,omitempty"`     /*备注*/
	IsSend      string       `json:"isSend,omitempty"`     /*变动是否推送给顾客*/
	CreatedBy   string       `json:"createdBy,omitempty"`
	UpdatedBy   string       `json:"updatedBy,omitempty"`
	CreatedAt   *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time   `json:"updatedAt,omitempty"`
	MileageDtls []MileageDtl `json:"mileageDtls,omitempty"`
}

type MileageDtl struct {
	Id         int64      `json:"id,omitempty"`        /*ID 主键*/
	MileageId  int64      `json:"mileageId,omitempty"` /*Mst表主键*/
	ItemId     int64      `json:"itemId,omitempty"`    /*详情Id*/
	PreItemId  int64      `json:"preItemId,omitempty"` /*原详情Id*/
	OfferNo    string     `json:"offerNo,omitempty"`   /*促销号*/
	Type       string     `json:"type,omitempty"`      /*类型*/
	SaleAmount float64    `json:"saleAmount"`          /*整单金额*/
	PayAmount  float64    `json:"payAmount"`           /*实付金额*/
	Point      float64    `json:"point"`               /*积分数量*/
	PointPrice float64    `json:"pointPrice"`          /*积分抵扣金额*/
	Remark     string     `json:"remark,omitempty"`    /*备注*/
	IsSend     string     `json:"isSend,omitempty"`    /*变动是否推送给顾客*/
	CreatedBy  string     `json:"createdBy,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedBy  string     `json:"updatedBy,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type MileageMall struct {
	TenantCode             string     `json:"tenantCode,omitempty"`   /*租户代码*/
	MallId                 int64      `json:"mallId,omitempty"`       /*购物中心代码*/
	MemberId               int64      `json:"memberId,omitempty"`     /*会员id*/
	CardTypeId             int64      `json:"cardTypeId"`             /*卡片类型*/
	GradeId                int64      `json:"gradeId"`                /*等级*/
	PreExitType            int64      `json:"PreExitType"`            /*预退出顾客类型，适用于eland集团5、6等级顾客*/
	Point                  float64    `json:"point"`                  /*总积分数量*/
	TotalAmount            float64    `json:"totalAmount"`            /*交易总金额*/
	TotalCount             int64      `json:"totalCount"`             /*消费总笔数*/
	RecentPurchaseDatetime *time.Time `json:"recentPurchaseDatetime"` /*最近消费时间*/
	OldMobile              string     `json:"oldMobile,omitempty"`
	CreatedBy              string     `json:"createdBy,omitempty"`
	UpdatedBy              string     `json:"updatedBy,omitempty"`
	CreatedAt              *time.Time `json:"createdAt,omitempty"`
	UpdatedAt              *time.Time `json:"updatedAt,omitempty"`
}

func (Mileage) GetMembershipMileages(ctx context.Context, tradeNo int64) ([]Mileage, error) {
	userClaim := auth.UserClaim{}.FromCtx(ctx)
	var resp struct {
		Result struct {
			Items      []Mileage `json:"items"`
			TotalCount int64     `json:"totalCount"`
		} `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/mileage?tradeNo=%v&tenantCode=%s",
		config.Config().Services.BenefitApi, tradeNo, userClaim.TenantCode)
	logrus.WithField("url", url).Info("url")
	time.Sleep(2 * time.Second)
	if err := RetryRestApi(ctx, resp, http.MethodGet, url, nil); err != nil {
		return nil, fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}

	return resp.Result.Items, nil
}

func (MileageMall) GetMembershipGrade(ctx context.Context, brandId, memberId int64, tenantCode string) (*MileageMall, error) {
	var resp struct {
		Result  *MileageMall `json:"result"`
		Success bool         `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/mileage-mall?tenantCode=%s&mallId=%v&memberId=%v",
		config.Config().Services.BenefitApi, tenantCode, brandId, memberId)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get mileage_mall")
		return nil, fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}

	return resp.Result, nil
}
