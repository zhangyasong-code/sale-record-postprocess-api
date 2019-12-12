package saleRecordFee

import (
	"context"
	"nhub/sale-record-postprocess-api/models"

	"github.com/sirupsen/logrus"
)

type SaleRecordFeeEventHandler struct {
}

func (h SaleRecordFeeEventHandler) Handle(ctx context.Context, a models.SaleRecordEvent) error {
	if err := setPostSaleRecordFee(ctx, a); err != nil {
		return err
	}
	return nil
}

func setPostSaleRecordFee(ctx context.Context, a models.SaleRecordEvent) error {
	postSaleRecordFees, err := PostSaleRecordFee{}.MakePostSaleRecordFeesEntity(ctx, a)
	if err != nil {
		return err
	}
	for _, postSaleRecordFee := range postSaleRecordFees {
		has, _, err := postSaleRecordFee.Get(ctx)
		if err != nil {
			logrus.WithField("Error", err).Info("GetpostSaleRecordFee error")
			return err
		}
		if has {
			continue
		} else {
			if err := postSaleRecordFee.Save(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
