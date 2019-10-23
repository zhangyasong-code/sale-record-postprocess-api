package sendCsl

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"nhub/sale-record-postprocess-api/models"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

type Send struct{}

type Payload struct {
	BrandCode   string `json:"brandCode"`
	ChannelType string `json:"channelType"`
	OrderId     int64  `json:"orderId"`
	RefundId    int64  `json:"refundId"`
}

func (Send) SendToCsl(ctx context.Context, event models.SaleRecordEvent) error {
	var resp struct {
		Success bool `json:"success"`
		Error   struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	payload, err := getPayload(event)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/v1/transaction/sale-csl", config.Config().Services.ClearanceAdapterForSaleRecordApi)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodPost, url, payload).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorMessage": resp.Error.Message,
		}).Error("Fail to SendToCsl")
		return fmt.Errorf("%s", resp.Error.Message)
	}

	return nil
}

func getPayload(event models.SaleRecordEvent) (Payload, error) {
	var payload Payload
	for i, dtl := range event.AssortedSaleRecordDtlList {
		if i == 0 {
			payload.BrandCode = dtl.BrandCode
		}
	}
	payload.ChannelType = event.TransactionChannelType
	payload.OrderId = event.OrderId
	payload.RefundId = event.RefundId
	return payload, nil
}
