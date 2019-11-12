package customer

import (
	"context"

	"nhub/sale-record-postprocess-api/models"

	"github.com/sirupsen/logrus"
)

type CustomerEventHandler struct {
}

func (h CustomerEventHandler) Handle(ctx context.Context, record models.SaleRecordEvent) error {
	if record.CustomerId == 0 || (record.Mileage == 0 && record.ObtainMileage == 0) {
		return nil
	}
	has, err := PostMileage{}.CheckOrderRefundExist(ctx, record.TransactionId)
	if err != nil {
		logrus.WithField("err", err).Info("CheckOrderRefundExist")
		return err
	}
	if has {
		logrus.WithFields(logrus.Fields{
			"err":           "TransactionId has exist.",
			"transactionId": record.TransactionId,
		})
		return nil
	}

	if err := saveMileageFromSaleRecord(ctx, record); err != nil {
		return err
	}
	return nil
}

func saveMileageFromSaleRecord(ctx context.Context, record models.SaleRecordEvent) error {
	mileages, err := Mileage{}.MakeMileage(ctx, record)
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
		}
	}
	return nil
}

func saveMileageFromMembership(ctx context.Context, record models.SaleRecordEvent) error {
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
		}
	}
	return nil
}
