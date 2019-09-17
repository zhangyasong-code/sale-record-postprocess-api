package saleRecordFee

import (
	"context"
	"nhub/sale-record-postprocess-api/customer"
	"nhub/sale-record-postprocess-api/models"
	"nhub/sale-record-postprocess-api/promotion"

	"github.com/sirupsen/logrus"

	"github.com/pangpanglabs/goutils/number"
)

func (PostSaleRecordFee) MakePostSaleRecordFeesEntity(ctx context.Context, a models.SaleRecordEvent) ([]PostSaleRecordFee, error) {
	var postSaleRecordFees []PostSaleRecordFee
	var eventFeeRate, appliedFeeRate, feeAmount, itemFeeRate float64
	var eventTypeCode string

	appliedFeeRate = 0
	eventTypeCode = ""
	eventFeeRate = 0
	for _, cartOffer := range a.CartOffers {
		if eventFeeRate != 0 {
			continue
		}
		promotionEvent, err := promotion.GetByNo(ctx, cartOffer.OfferNo)
		if err != nil {
			logrus.WithField("Error", err).Info("GetPromotionEvent error")
			return nil, err
		}
		eventTypeCode = promotionEvent.EventTypeCode
		if eventTypeCode == "01" || eventTypeCode == "02" || eventTypeCode == "03" {
			eventFeeRate = promotionEvent.FeeRate
			if eventFeeRate <= 0 {
				postFailCreateSaleFee := &PostFailCreateSaleFee{TransactionId: a.TransactionId, IsProcessed: false}
				has, _, err := postFailCreateSaleFee.Get(ctx)
				if err != nil {
					return nil, err
				}
				if !has {
					if err := postFailCreateSaleFee.Save(ctx); err != nil {
						return nil, err
					}
				}
				return nil, nil
			}
		}
	}

	for _, assortedSaleRecordDtl := range a.AssortedSaleRecordDtlList {
		feeAmount = 0
		// Use the offerNo to query promotionEvent
		if eventFeeRate == 0 && len(assortedSaleRecordDtl.ItemOffers) != 0 {
			for _, ItemOffer := range assortedSaleRecordDtl.ItemOffers {
				if eventFeeRate != 0 {
					continue
				}
				promotionEvent, err := promotion.GetByNo(ctx, ItemOffer.OfferNo)
				if err != nil {
					logrus.WithField("Error", err).Info("GetPromotionEvent error")
					return nil, err
				}
				eventTypeCode = promotionEvent.EventTypeCode
				if eventTypeCode == "01" || eventTypeCode == "02" || eventTypeCode == "03" {
					eventFeeRate = promotionEvent.FeeRate
					if eventFeeRate <= 0 {
						postFailCreateSaleFee := &PostFailCreateSaleFee{TransactionId: a.TransactionId, IsProcessed: false}
						has, _, err := postFailCreateSaleFee.Get(ctx)
						if err != nil {
							return nil, err
						}
						if !has {
							if err := postFailCreateSaleFee.Save(ctx); err != nil {
								return nil, err
							}
						}
						return nil, nil
					}
				}
			}
		}
		appliedFeeRate = eventFeeRate
		itemFeeRate = assortedSaleRecordDtl.FeeRate
		// eventFeeRate 优先级大于 itemFeeRate
		if appliedFeeRate == 0 && itemFeeRate > 0 {
			appliedFeeRate = itemFeeRate
		}
		useType := customer.UseTypeUsed
		if assortedSaleRecordDtl.RefundItemId != 0 {
			useType = customer.UseTypeUsedCancel
		}
		// Use the OrderItemId to query Mileage and MileagePrice
		_, mileagePrice, err := customer.PostMileageDtl{}.GetPostMileageDtl(ctx, 0, assortedSaleRecordDtl.OrderItemId, assortedSaleRecordDtl.RefundItemId, useType)
		if err != nil {
			logrus.WithFields(logrus.Fields{"OrderItemId": assortedSaleRecordDtl.OrderItemId, "RefundItemId": assortedSaleRecordDtl.RefundItemId, "Error": err}).Error("GetOrgMileageDtl failed!")
			return nil, err
		}

		// ((TotalDistributedPaymentPrice - mileagePrice) * appliedFeeRate) / 100
		feeAmount = number.ToFixed(((assortedSaleRecordDtl.DistributedPrice.TotalDistributedPaymentPrice-mileagePrice.PointPrice)*appliedFeeRate)/100, nil)
		postSaleRecordFees = append(
			postSaleRecordFees,
			PostSaleRecordFee{
				TransactionId:          a.TransactionId,
				TransactionDtlId:       assortedSaleRecordDtl.Id,
				OrderId:                a.OrderId,
				OrderItemId:            assortedSaleRecordDtl.OrderItemId,
				RefundId:               a.RefundId,
				RefundItemId:           assortedSaleRecordDtl.RefundItemId,
				CustomerId:             a.CustomerId,
				StoreId:                a.StoreId,
				TotalSalePrice:         assortedSaleRecordDtl.TotalPrice.SalePrice,
				TotalPaymentPrice:      assortedSaleRecordDtl.DistributedPrice.TotalDistributedPaymentPrice,
				Mileage:                mileagePrice.Point,
				MileagePrice:           mileagePrice.PointPrice,
				ItemFeeRate:            itemFeeRate,
				ItemFee:                assortedSaleRecordDtl.ItemFee,
				EventFeeRate:           eventFeeRate,
				AppliedFeeRate:         appliedFeeRate,
				FeeAmount:              feeAmount,
				TransactionChannelType: a.TransactionChannelType,
			},
		)
	}
	return postSaleRecordFees, nil
}
