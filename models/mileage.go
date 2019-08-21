package models

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

type Mileage struct {
	Id         int64      `json:"id,omitempty"`         /*ID 主键*/
	TenantCode string     `json:"tenantCode,omitempty"` /*租户代码*/
	MallId     string     `json:"mallId,omitempty"`     /*购物中心代码*/
	StoreId    string     `json:"storeId,omitempty"`    /*店铺代码*/
	MemberId   int64      `json:"memberId"`             /*会员id*/
	Type       string     `json:"type,omitempty"`       /*类型*/
	ChannelId  int64      `json:"channelId"`            /*渠道*/
	TradeDate  time.Time  `json:"tradeDate,omitempty"`  /*交易日期*/
	TradeNo    string     `json:"tradeNo,omitempty"`    /*交易单号*/
	PreTradeNo string     `json:"preTradeNo,omitempty"` /*原交易单号*/
	SaleAmount float64    `json:"saleAmount"`           /*整单金额*/
	PayAmount  float64    `json:"payAmount"`            /*实付金额*/
	Point      float64    `json:"point"`                /*积分数量*/
	Remark     string     `json:"remark,omitempty"`     /*备注*/
	IsSend     string     `json:"isSend,omitempty"`     /*变动是否推送给顾客*/
	CreatedBy  string     `json:"createdBy,omitempty"`
	UpdatedBy  string     `json:"updatedBy,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

func (Mileage) GetMembershipMileages(ctx context.Context, tradeNo int64, mileageType string) ([]Mileage, error) {
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
	url := fmt.Sprintf("%s/v1/mileage?tradeNo=%v&tenantCode=%s&type=%s",
		config.Config().Services.MembershipApi, tradeNo, userClaim.TenantCode, mileageType)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get mileages")
		return nil, fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}

	return resp.Result.Items, nil
}
