package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/pangpanglabs/goutils/number"
)

func (PostSaleRecordFee) MakePostSaleRecordFeesEntiiy(ctx context.Context, a SaleRecordEvent) ([]PostSaleRecordFee, error) {
	var postSaleRecordFees []PostSaleRecordFee
	var eventFeeRate, appliedFeeRate, feeAmount float64
	var eventTypeCode string
	for _, assortedSaleRecordDtl := range a.AssortedSaleRecordDtlList {
		eventFeeRate = 0
		appliedFeeRate = 0
		feeAmount = 0
		eventTypeCode = ""
		// Use the offerNo to query promotionEvent
		if len(assortedSaleRecordDtl.ItemOffers) != 0 {
			for _, ItemOffer := range assortedSaleRecordDtl.ItemOffers {
				promotionEvent, err := PostSaleRecordFee{}.GetPromotionEvent(ctx, ItemOffer.OfferNo)
				if err != nil {
					logrus.WithField("Error", err).Info("GetPromotionEvent error")
					return nil, err
				}
				eventTypeCode = promotionEvent.EventTypeCode
				if eventTypeCode == "01" || eventTypeCode == "02" || eventTypeCode == "03" {
					eventFeeRate += promotionEvent.FeeRate
					if eventFeeRate <= 0 {
						postFailCreateSaleFee := &PostFailCreateSaleFee{TransactionId: a.TransactionId, IsProcessed: false}
						if err := postFailCreateSaleFee.Save(ctx); err != nil {
							return nil, err
						}
					}
				}
			}
			appliedFeeRate = eventFeeRate
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
			logrus.WithField("TransactionId", a.TransactionId).Info("Add FailCreateSaleFee data")
			postFailCreateSaleFee := &PostFailCreateSaleFee{TransactionId: a.TransactionId, IsProcessed: false}
			if err := postFailCreateSaleFee.Save(ctx); err != nil {
				return nil, err
			}
			return nil, nil
		}

		// Use the OrderItemId to query Mileage and MileagePrice
		mileagePrice, err := PostSaleRecordFee{}.GetPostMileageDtl(ctx, assortedSaleRecordDtl.OrderItemId, assortedSaleRecordDtl.RefundItemId)
		if err != nil {
			logrus.WithFields(logrus.Fields{"OrderItemId": assortedSaleRecordDtl.OrderItemId, "RefundItemId": assortedSaleRecordDtl.RefundItemId, "Error": err}).Error("GetOrgMileageDtl failed!")
			return nil, err
		}

		// ((TotalDistributedPaymentPrice - mileagePrice) * appliedFeeRate) / 100
		feeAmount = number.ToFixed(((assortedSaleRecordDtl.DistributedPrice.TotalDistributedPaymentPrice-mileagePrice.PointAmount)*appliedFeeRate)/100, nil)
		postSaleRecordFees = append(
			postSaleRecordFees,
			PostSaleRecordFee{
				TransactionId:          a.TransactionId,
				SaleRecordDtlId:        assortedSaleRecordDtl.Id,
				OrderId:                a.OrderId,
				OrderItemId:            assortedSaleRecordDtl.OrderItemId,
				RefundId:               a.RefundId,
				RefundItemId:           assortedSaleRecordDtl.RefundItemId,
				CustomerId:             a.CustomerId,
				StoreId:                a.StoreId,
				TotalSalePrice:         assortedSaleRecordDtl.TotalPrice.SalePrice,
				TotalPaymentPrice:      assortedSaleRecordDtl.DistributedPrice.TotalDistributedPaymentPrice,
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
