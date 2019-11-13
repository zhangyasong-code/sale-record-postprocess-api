package customer

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"time"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

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
