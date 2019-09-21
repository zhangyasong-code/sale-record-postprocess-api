package payamt

import (
	"context"
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
	pays, err := Pay{}.GetPayamt(ctx, record.OrderId)
	if len(pays) > 0 {
		for _, pay := range pays {
			paymentCode := "11"
			creditCardFirmCode := ""
			if pay.PayMethod == "CASH" {
				paymentCode = "11"
			} else if pay.PayMethod == "WXPAY" {
				paymentCode = "01"
			} else if pay.PayMethod == "wechat.prepay" {
				paymentCode = "01"
			} else if pay.PayMethod == "ALIPAY" {
				paymentCode = "02"
			} else if pay.PayMethod == "CREDITCARD" {
				paymentCode = "23"
				creditCardFirmCode = "01"
			}
			postPayment = append(postPayment, PostPayment{
				TransactionId:      record.TransactionId,
				SeqNo:              pay.SeqNo,
				PaymentCode:        paymentCode,
				PaymentAmt:         pay.PayAmt,
				InUserID:           pay.CreatedBy,
				InDateTime:         pay.CreatedAt,
				ModiUserID:         pay.CreatedBy,
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
