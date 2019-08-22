package models

import (
	"context"

	"github.com/pangpanglabs/goutils/number"
)

type CustomerEventHandler struct {
}

func (h CustomerEventHandler) Handle(ctx context.Context, a SaleRecordEvent) error {
	if err := setAccumulateMileage(ctx, a); err != nil {
		return err
	}
	if a.Mileage != 0 {
		if err := setUsedMileage(ctx, a); err != nil {
			return err
		}
	}
	if err := setPostSaleRecordFee(ctx, a); err != nil {
		return err
	}

	return nil
}

func setAccumulateMileage(ctx context.Context, a SaleRecordEvent) error {
	tradeNo := a.OrderId
	if a.RefundId != 0 {
		tradeNo = a.RefundId
	}

	accumulateMileages, err := Mileage{}.GetMembershipMileages(ctx, tradeNo, "E")
	if err != nil {
		return err
	}
	for _, accumulateMileage := range accumulateMileages {
		if accumulateMileage.Point != 0 {
			orgMileage := OrgMileage{}.MakeOrgMileage(&accumulateMileage, a)
			if err := orgMileage.Create(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func setUsedMileage(ctx context.Context, a SaleRecordEvent) error {
	orgMileage := OrgMileage{}.MakeOrgMileage(nil, a)
	if err := orgMileage.Create(ctx); err != nil {
		return err
	}
	orgMileageDtls := make([]OrgMileageDtl, 0)
	var lastPoint, lastPointAmount float64
	for i, recordDtl := range a.AssortedSaleRecordDtls {
		var currentPoint, currentPointAmount float64
		ratio := recordDtl.TotalSalePrice / a.TotalPrice.SalePrice
		if i == len(a.AssortedSaleRecordDtls) {
			currentPoint = lastPoint
			currentPointAmount = lastPointAmount
		} else {
			currentPoint = number.ToFixed(a.Mileage*ratio, nil)
			currentPointAmount = number.ToFixed(a.MileagePrice*ratio, nil)
			lastPoint = a.Mileage - currentPoint
			lastPointAmount = a.MileagePrice - currentPointAmount
		}
		orgMileageDtl := OrgMileageDtl{}.MakeOrgMileageDtl(orgMileage.Id,
			orgMileage.MileageType, currentPoint, currentPointAmount, recordDtl)
		orgMileageDtls = append(orgMileageDtls, *orgMileageDtl)
	}
	if err := (OrgMileageDtl{}).CreateBatch(ctx, orgMileageDtls); err != nil {
		return err
	}
	return nil
}

func setPostSaleRecordFee(ctx context.Context, a SaleRecordEvent) error {
	postSaleRecordFees, err := PostSaleRecordFee{}.MakePostSaleRecordFeesEntiiy(ctx, a)
	if err != nil {
		return err
	}
	for _, postSaleRecordFee := range postSaleRecordFees {
		if err := postSaleRecordFee.Save(ctx); err != nil {
			return err
		}
	}
	return nil
}
