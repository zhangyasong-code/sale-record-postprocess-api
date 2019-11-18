package customer

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

type BrandMember struct {
	BrandId       int64 `json:"brandId,omitempty"`       /*品牌Id*/
	MemberBrandId int64 `json:"memberBrandId,omitempty"` /*会员品牌Id*/
}

func (BrandMember) GetBrandMembers(ctx context.Context, brandIds string) ([]BrandMember, error) {
	var resp struct {
		Result  []BrandMember `json:"result"`
		Success bool          `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/brand/member?brandIds=%s", config.Config().Services.MemberApi, brandIds)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get brandmembers")
		return nil, fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}

	return resp.Result, nil
}
