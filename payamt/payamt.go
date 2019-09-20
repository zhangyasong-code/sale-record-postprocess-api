package payamt

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

type Pay struct {
	Id              int64     `json:"id" xorm:"int64 notnull autoincr pk 'id'"`
	OrderId         int64     `json:"orderId"`
	CartId          int64     `json:"cartId"`
	SeqNo           int64     `json:"seqNo"`
	PayMethod       string    `json:"payMethod"`
	EId             int64     `json:"eId" xorm:"int64 'e_id'"`
	PayAmt          float64   `json:"payAmt"`
	PayInfo         string    `json:"payInfo" xorm:"varchar(4000)"`
	Status          string    `json:"status"`
	TenantCode      string    `json:"tenantCode"`
	ShopCode        string    `json:"shopCode"`
	RefundOrderId   int64     `json:"refundOrderId"`
	RefundAmt       float64   `json:"refundAmt"`
	CreatedAt       time.Time `json:"createdAt" xorm:"created"`
	CreatedBy       string    `json:"createdBy"`
	Source          string    `json:"source"`
	OutTradeNo      string    `json:"outTradeNo"`
	VoucherNo       string    `json:"voucherNo"`
	ReferenceNo     string    `json:"referenceNo"`
	TranceNo        string    `json:"tranceNo"`
	Terminal        string    `json:"terminal"`
	MerchantId      string    `json:"merchantId"`   //银行签约商户号
	MerchantName    string    `json:"merchantName"` //银行签约商户名称
	CardNo          string    `json:"cardNo"`
	CardType        string    `json:"cardType"`
	OriginalId      int64     `json:"originalId"`
	Uid             string    `json:"uid"`
	MemberAmt       float64   `json:"memberAmt"`
	OutRefundNo     string    `json:"outRefundNo"`
	MdiscountAmount float64   `json:"mdiscountAmount"`
}

func (Pay) GetPayamt(ctx context.Context, orderId int64) ([]Pay, error) {
	var resp struct {
		Result  []Pay `json:"result"`
		Success bool  `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/query/orderId/%v", config.Config().Services.PayamtApi, orderId)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get payamt")
		return nil, fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}

	return resp.Result, nil
}
