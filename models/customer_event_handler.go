package models

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CustomerEventHandler struct {
}

func (h CustomerEventHandler) Handle(ctx context.Context, record SaleRecordEvent) error {
	has, err := PostMileage{}.CheckOrderRefundExist(ctx, record.TransactionId)
	if err != nil {
		logrus.WithField("err", err).Info("CheckOrderRefundExist")
		return err
	}
	if has {
		return fmt.Errorf("TransactionId(%v) has exist.", record.TransactionId)
	}

	tradeNo := record.OrderId
	if record.RefundId != 0 {
		tradeNo = record.RefundId
	}
	mileages, err := Mileage{}.GetMembershipMileages(ctx, tradeNo)
	if err != nil {
		return err
	}
	for _, mileage := range mileages {
		if mileage.Point != 0 {
			postMileage, err := PostMileage{}.MakePostMileage(ctx, mileage, record)
			if err != nil {
				return err
			}
			if err := postMileage.Create(ctx); err != nil {
				return err
			}
			postMileageDtls := PostMileageDtl{}.MakePostMileageDtls(postMileage, mileage.MileageDtls, record.AssortedSaleRecordDtlList)
			if err := (PostMileageDtl{}).CreateBatch(ctx, postMileageDtls); err != nil {
				return err
			}
		}
	}
	return nil
}
