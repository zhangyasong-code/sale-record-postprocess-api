package payamt

import (
	"context"
	"math"
	"nhub/sale-record-postprocess-api/models"
)

type PayAmtEventHandler struct {
}

func (h PayAmtEventHandler) Handle(ctx context.Context, record models.SaleRecordEvent) error {
	has, err := (PostPayment{}).checkExist(ctx, record.TransactionId)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	var postPayment []PostPayment
	if len(record.Payments) > 0 {
		for _, pay := range record.Payments {
			if pay.PayMethod == "MILEAGE" {
				continue
			}
			paymentCode := "11"
			creditCardFirmCode := ""
			if pay.PayMethod == "CASH" {
				paymentCode = "11"
			} else if pay.PayMethod == "WXPAY" {
				paymentCode = "O1"
			} else if pay.PayMethod == "wechat.prepay" {
				paymentCode = "O1"
			} else if pay.PayMethod == "ALIPAY" {
				paymentCode = "O2"
			} else if pay.PayMethod == "CREDITCARD" {
				paymentCode = "12"
				creditCardFirmCode = "01"
			}
			postPayment = append(postPayment, PostPayment{
				TransactionId:      record.TransactionId,
				SeqNo:              pay.SeqNo,
				PaymentCode:        paymentCode,
				PaymentAmt:         math.Abs(pay.PayAmt),
				InDateTime:         pay.CreatedAt,
				ModiDateTime:       pay.CreatedAt,
				CreditCardFirmCode: creditCardFirmCode,
			})
		}
	}

	//2.保存到数据库
	if err = (PostPayment{}).createInArrary(ctx, postPayment); err != nil {
		return err
	}

	return nil
}
