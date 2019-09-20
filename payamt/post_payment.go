package payamt

import (
	"context"
	"nhub/sale-record-postprocess-api/factory"
	"time"
)

type PostPayment struct {
	Id                 int64     `json:"id" xorm:"pk"`
	TransactionId      int64     `json:"transactionId"`
	SeqNo              int64     `json:"seqNo"`
	PaymentCode        string    `json:"paymentCode"`
	PaymentAmt         float64   `json:"paymentAmt"`
	InUserID           string    `json:"inUserId"`
	InDateTime         time.Time `json:"inDateTime" xorm:"created"`
	ModiUserID         string    `json:"modiUserID"`
	ModiDateTime       time.Time `json:"modiDateTime" xorm:"created"`
	CreditCardFirmCode string    `json:"creditCardFirmCode"`
}

func (PostPayment) createInArrary(ctx context.Context, postPayment []PostPayment) error {
	if _, err := factory.SaleRecordDB(ctx).Insert(&postPayment); err != nil {
		return err
	}
	return nil
}
