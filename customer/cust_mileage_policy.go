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

type CustMileagePolicy struct {
	BrandCode                            string  `json:"brandCode"`
	CustMileagePolicyNo                  int64   `json:"custMileagePolicyNo"`
	MileageUseChk                        bool    `json:"mileageUseChk"`
	UsePermitChk                         bool    `json:"usePermitChk"`
	MinUseAmt                            float64 `json:"minUseAmt"`
	MaxUseRate                           float64 `json:"maxUseRate"`
	PurchaseStartDate                    string  `json:"purchaseStartDate"`
	PurchaseEndDate                      string  `json:"purchaseEndDate"`
	UseMileageOperationUnitCode          int64   `json:"useMileageOperationUnitCode"`
	AccumulationMileageOperationUnitCode int64   `json:"accumulationMileageOperationUnitCode"`
}

func (CustMileagePolicy) GetCustMileagePolicy(ctx context.Context, brandCode string) (*CustMileagePolicy, error) {
	var resp struct {
		Result  *CustMileagePolicy `json:"result"`
		Success bool               `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/custmileagepolicy/%s", config.Config().Services.CslCustomerApi, brandCode)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get custmileagepolicy")
		return nil, fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}

	return resp.Result, nil
}
