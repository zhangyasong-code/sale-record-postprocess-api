package refundApproval

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"time"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
)

type RefundApprovalInput struct {
	TenantCode string    `json:"tenantCode"`
	StoreId    int64     `json:"storeId"`
	OrderId    int64     `json:"orderId"`
	RefundId   int64     `json:"refundId"`
	RefundTime time.Time `json:"refundTime"`
}

func Check(ctx context.Context, tenantCode string, storeId, orderId, refundId int64, refundTime time.Time) (bool, error) {
	var resp struct {
		Result struct {
			Status int `json:"status"`
		} `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int64       `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}
	refundApprovalInput := RefundApprovalInput{
		TenantCode: tenantCode,
		StoreId:    storeId,
		OrderId:    orderId,
		RefundId:   refundId,
		RefundTime: refundTime,
	}

	url := fmt.Sprintf("%s/v1/check", config.Config().Services.RefundApprovalApi)
	if _, err := httpreq.New(http.MethodPost, url, refundApprovalInput).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return false, err
	}

	if resp.Success {
		return true, nil
	}
	return false, nil
}
