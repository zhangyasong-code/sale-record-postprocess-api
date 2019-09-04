package models

import (
	"context"

	"github.com/pangpanglabs/goutils/number"
	"github.com/sirupsen/logrus"
)

type CustomerEventHandler struct {
}

func (h CustomerEventHandler) Handle(ctx context.Context, a SaleRecordEvent) error {
	//mileage
	if err := setPostMileage(ctx, a); err != nil {
		return err
	}
	return nil
}

func setPostMileage(ctx context.Context, a SaleRecordEvent) error {
	has, err := PostMileage{}.CheckOrderRefundExist(ctx, a.TransactionId)
	if err != nil {
		logrus.WithField("err", err).Info("CheckOrderRefundExist")
		return err
	}
	if has {
		return nil
	}

	if err := setAccumulateMileage(ctx, a); err != nil {
		logrus.WithField("err", err).Info("setAccumulateMileage")
		return err
	}
	if a.Mileage != 0 {
		if err := setUsedMileage(ctx, a); err != nil {
			logrus.WithField("err", err).Info("setUsedMileage")
			return err
		}
	}
	return nil
}

func setAccumulateMileage(ctx context.Context, saleRecordEvent SaleRecordEvent) error {
	tradeNo := saleRecordEvent.OrderId
	if saleRecordEvent.RefundId != 0 {
		tradeNo = saleRecordEvent.RefundId
	}

	accumulateMileages, err := Mileage{}.GetMembershipMileages(ctx, tradeNo, "E")
	if err != nil {
		return err
	}
	for _, accumulateMileage := range accumulateMileages {
		if accumulateMileage.Point != 0 {
			mileageMall, err := MileageMall{}.GetMembershipGrade(ctx, accumulateMileage.MallId, accumulateMileage.MemberId, accumulateMileage.TenantCode)
			if err != nil {
				return err
			}
			postMileage := PostMileage{}.MakePostMileage(&accumulateMileage, mileageMall, saleRecordEvent)
			if err := postMileage.Create(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func setUsedMileage(ctx context.Context, saleRecordEvent SaleRecordEvent) error {
	postMileage := PostMileage{}.MakePostMileage(nil, nil, saleRecordEvent)
	if err := postMileage.Create(ctx); err != nil {
		return err
	}
	postMileageDtls := make([]PostMileageDtl, 0)
	var lastPoint, lastPointAmount float64
	for i, recordDtl := range saleRecordEvent.AssortedSaleRecordDtlList {
		var currentPoint, currentPointAmount float64
		ratio := recordDtl.TotalPrice.SalePrice / saleRecordEvent.TotalPrice.SalePrice
		if i == len(saleRecordEvent.AssortedSaleRecordDtlList) {
			currentPoint = lastPoint
			currentPointAmount = lastPointAmount
		} else {
			currentPoint = number.ToFixed(saleRecordEvent.Mileage*ratio, nil)
			currentPointAmount = number.ToFixed(saleRecordEvent.MileagePrice*ratio, nil)
			lastPoint = saleRecordEvent.Mileage - currentPoint
			lastPointAmount = saleRecordEvent.MileagePrice - currentPointAmount
		}
		postMileageDtl := PostMileageDtl{}.MakePostMileageDtl(postMileage.Id,
			postMileage.UseType, currentPoint, currentPointAmount, recordDtl)
		postMileageDtls = append(postMileageDtls, *postMileageDtl)
	}
	if err := (PostMileageDtl{}).CreateBatch(ctx, postMileageDtls); err != nil {
		return err
	}
	return nil
}
