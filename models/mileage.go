package models

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"nomni/utils/auth"

	membership "membership-api/models"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

func GetMembershipMileages(ctx context.Context, tradeNo int64, mileageType string) ([]membership.Mileage, error) {
	userClaim := auth.UserClaim{}.FromCtx(ctx)
	var resp struct {
		Result struct {
			Items      []membership.Mileage `json:"items"`
			TotalCount int64                `json:"totalCount"`
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
