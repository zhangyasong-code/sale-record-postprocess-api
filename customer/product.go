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

type BrandInfo struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (BrandInfo) GetBrandInfo(ctx context.Context, brandId int64) (string, error) {
	var resp struct {
		Result  *BrandInfo `json:"result"`
		Success bool       `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}

	url := fmt.Sprintf("%s/v1/brands/%v", config.Config().Services.ProductApi, brandId)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return "", err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get brandcode")
		return "", fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}
	return resp.Result.Code, nil
}
