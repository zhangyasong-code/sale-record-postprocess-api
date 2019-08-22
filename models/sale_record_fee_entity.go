package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/pangpanglabs/goutils/number"
)

func (PostSaleRecordFee) MakePostSaleRecordFeesEntiiy(ctx context.Context, a SaleRecordEvent) ([]PostSaleRecordFee, error) {
	var postSaleRecordFees []PostSaleRecordFee
	for _, assortedSaleRecordDtl := range a.AssortedSaleRecordDtls {
		var eventFeeRate, appliedFeeRate, feeAmount float64
		var eventType string

		// Use the offerNo to query promotionEvent
		if assortedSaleRecordDtl.Offer.OfferId != 0 {
			promotionEvent, err := PostSaleRecordFee{}.GetPromotionEvent(ctx, assortedSaleRecordDtl.Offer.OfferNo)
			if err != nil {
				logrus.WithField("Error", err).Info("GetPromotionEvent error")
				return nil, err
			}
			eventFeeRate = promotionEvent.FeeRate
			eventType = promotionEvent.EventTypeCode
			if eventType != "" && eventFeeRate > 0 {
				appliedFeeRate = promotionEvent.FeeRate
			}
		}

		// Use the brandId and TransactionCreateDate to get FeeRate in contracts
		contractFeeRate, err := PostSaleRecordFee{}.GetContractFeeRate(ctx, a.StoreId, assortedSaleRecordDtl.BrandId, a.TransactionCreateDate)
		if err != nil {
			return nil, err
		}
		// eventFeeRate 优先级大于 contractFeeRate
		if appliedFeeRate == 0 && contractFeeRate > 0 {
			appliedFeeRate = contractFeeRate
		} else if appliedFeeRate == 0 && contractFeeRate == 0 {
			// Add one Case when eventFeeRate and contractFeeRate is 0
		}

		// Use the OrderItemId to query Mileage and MileagePrice
		mileagePrice, err := PostSaleRecordFee{}.GetOrgMileageDtl(ctx, assortedSaleRecordDtl.OrderItemId, assortedSaleRecordDtl.RefundItemId)
		if err != nil {
			logrus.WithFields(logrus.Fields{"OrderItemId": assortedSaleRecordDtl.OrderItemId, "RefundItemId": assortedSaleRecordDtl.RefundItemId, "Error": err}).Error("GetOrgMileageDtl failed!")
			return nil, err
		}

		// ((totalSalePrice - mileagePrice) * appliedFeeRate) / 100
		feeAmount = number.ToFixed(((assortedSaleRecordDtl.TotalSalePrice-mileagePrice.PointAmount)*appliedFeeRate)/100, nil)
		postSaleRecordFees = append(
			postSaleRecordFees,
			PostSaleRecordFee{
				TransactionId:          a.TransactionId,
				SaleRecordDtlId:        assortedSaleRecordDtl.Id,
				SaleRecordOfferId:      assortedSaleRecordDtl.Offer.OfferId,
				OrderId:                a.OrderId,
				OrderItemId:            assortedSaleRecordDtl.OrderItemId,
				RefundId:               a.RefundId,
				RefundItemId:           assortedSaleRecordDtl.RefundItemId,
				CustomerId:             a.CustomerId,
				StoreId:                a.StoreId,
				EventType:              eventType,
				TotalSalePrice:         assortedSaleRecordDtl.TotalSalePrice,
				TotalPaymentPrice:      assortedSaleRecordDtl.TotalDistributedPaymentPrice,
				Mileage:                mileagePrice.Point,
				MileagePrice:           mileagePrice.PointAmount,
				ContractFeeRate:        contractFeeRate,
				EventFeeRate:           eventFeeRate,
				AppliedFeeRate:         appliedFeeRate,
				FeeAmount:              feeAmount,
				TransactionChannelType: a.TransactionChannelType,
			},
		)
	}
	return postSaleRecordFees, nil
}
