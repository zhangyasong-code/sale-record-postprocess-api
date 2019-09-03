package models

import (
	"context"
	"math"
	"nhub/sale-record-postprocess-api/promotion"
)

const (
	SaleEventGive          string  = "01" //送
	SaleEventReduce        string  = "02" //减
	SaleEventDiscount      string  = "03" //折
	CustEventBrandDiscount string  = "B"  //品牌折扣型
	CustEventCoupon        string  = "C"  //打折券event
	CustEventGift          string  = "G"  //赠品型
	CustEventMileage       string  = "M"  //积分型
	CustEventCustmoer      string  = "P"  //优秀顾客型
	CustEventAmountGift    string  = "R"  //赠品型-按购买金额
	CustEventVip           string  = "V"  //百货店vip型
	DiscountRate           float64 = 0.87 //是否算折扣金额的扣率
)

type SalesPersonEventHandler struct {
}

func (h SalesPersonEventHandler) Handle(ctx context.Context, s SaleRecordEvent) error {

	for i := 0; i < len(s.AssortedSaleRecordDtlList); i++ {
		saleAmountDtl := SaleRecordDtlSalesmanAmount{
			TransactionId:               s.TransactionId,
			SaleRecordDtlId:             s.AssortedSaleRecordDtlList[i].Id,
			OrderId:                     s.OrderId,
			RefundId:                    s.RefundId,
			StoreId:                     s.StoreId,
			SalesmanId:                  s.SalesmanId,
			TotalListPrice:              s.AssortedSaleRecordDtlList[i].TotalPrice.ListPrice,
			TotalSalePrice:              s.AssortedSaleRecordDtlList[i].TotalPrice.SalePrice,
			TotalDiscountPrice:          s.AssortedSaleRecordDtlList[i].TotalPrice.DiscountPrice,
			TotalDiscountCartOfferPrice: s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedCartOfferPrice,
			TotalPaymentPrice:           s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedPaymentPrice,
			TransactionChannelType:      s.TransactionChannelType,
			SalesmanSaleDiscountRate:    0,
			SalesmanSaleAmount:          s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedPaymentPrice,
		}
		//查询使用积分
		has, e, err := PostMileageDtl{}.GetByKey(ctx, s.AssortedSaleRecordDtlList[i].Id, s.AssortedSaleRecordDtlList[i].OrderItemId, s.AssortedSaleRecordDtlList[i].RefundItemId)
		if err != nil {
			return err
		}
		if has {
			saleAmountDtl.Mileage = e.Point
			saleAmountDtl.MileagePrice = e.PointAmount
		} else {
			saleAmountDtl.Mileage = 0
			saleAmountDtl.MileagePrice = 0
		}
		//计算营业员业绩金额-SalesmanSaleAmount
		offers := []SaleRecordDtlOffer{}
		if s.TotalPrice.DiscountPrice == 0 {
			saleAmountDtl.SalesmanSaleAmount = saleAmountDtl.TotalPaymentPrice
			saleAmountDtl.SalesmanNormalSaleAmount = saleAmountDtl.TotalPaymentPrice
			saleAmountDtl.SalesmanDiscountSaleAmount = 0
		} else {
			//查询使用优惠的类型
			for n := 0; n < len(s.AssortedSaleRecordDtlList[i].ItemOffers); n++ {
				p, err := promotion.GetByNo(ctx, s.AssortedSaleRecordDtlList[i].ItemOffers[n].OfferNo)
				if err != nil {
					return err
				}
				offer := SaleRecordDtlOffer{
					OfferId:         s.AssortedSaleRecordDtlList[i].ItemOffers[n].OfferId,
					EventType:       p.EventTypeCode,
					SaleBaseAmt:     p.SaleBaseAmt,
					DiscountBaseAmt: p.DiscountBaseAmt,
					DiscountRate:    p.DiscountRate,
				}
				offers = append(offers, offer)
				//计算实际销售额活动扣率
				if p.EventTypeCode == SaleEventGive { //送活动
					saleAmountDtl.SalesmanSaleDiscountRate = math.Floor((p.DiscountBaseAmt/(p.SaleBaseAmt*0.1+p.SaleBaseAmt+p.DiscountBaseAmt))*100) / 100
				}
				switch p.EventTypeCode {
				case SaleEventReduce, SaleEventDiscount: //减 折
					saleAmountDtl.SalesmanSaleAmount = saleAmountDtl.TotalPaymentPrice - saleAmountDtl.MileagePrice
				case SaleEventGive: //送
					saleAmountDtl.SalesmanSaleAmount = math.Floor((saleAmountDtl.TotalPaymentPrice-saleAmountDtl.MileagePrice-saleAmountDtl.TotalPaymentPrice*saleAmountDtl.SalesmanSaleDiscountRate)*100) / 100
				}
			}
			//拆分业绩金额-正常业绩和折扣业绩
			var saleEventTypeCode, primaryEventTypeCode, secondaryEventTypeCode string
			for d := 0; d < len(offers); d++ {
				switch offers[d].EventType {
				case SaleEventGive, SaleEventReduce, SaleEventDiscount: //活动-送  减  折
					saleEventTypeCode = offers[d].EventType
				case CustEventBrandDiscount, CustEventCoupon, CustEventCustmoer, CustEventVip: //基本顾客event-品牌折扣型 打折券 优秀顾客型 百货店VIP
					primaryEventTypeCode = offers[d].EventType
				case CustEventGift, CustEventMileage, CustEventAmountGift: //附加顾客event-赠品型 积分型 赠品型（按购买金额）
					secondaryEventTypeCode = offers[d].EventType
				}
			}
			if isNormalAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode) {
				saleAmountDtl.SalesmanNormalSaleAmount = saleAmountDtl.SalesmanSaleAmount
				saleAmountDtl.SalesmanDiscountSaleAmount = 0
			} else if isDiscountAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode) {
				saleAmountDtl.SalesmanNormalSaleAmount = 0
				saleAmountDtl.SalesmanDiscountSaleAmount = saleAmountDtl.SalesmanSaleAmount
			}
		}
		if err := (&saleAmountDtl).Create(ctx); err != nil {
			return err
		}
		for e := 0; e < len(offers); e++ {
			offers[e].SalesmanAmountId = saleAmountDtl.Id
			if err := (&offers[e]).Create(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func isNormalAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode string) bool {
	return (primaryEventTypeCode == CustEventBrandDiscount || primaryEventTypeCode == CustEventCoupon || primaryEventTypeCode == CustEventCustmoer) ||
		((saleEventTypeCode == SaleEventDiscount || primaryEventTypeCode == CustEventVip) && (saleAmountDtl.SalesmanSaleAmount/saleAmountDtl.TotalListPrice) > DiscountRate) ||
		(saleEventTypeCode == "" && primaryEventTypeCode == "" && (secondaryEventTypeCode == CustEventMileage || secondaryEventTypeCode == ""))
}
func isDiscountAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode string) bool {
	return ((saleEventTypeCode == SaleEventGive || saleEventTypeCode == SaleEventReduce) || (secondaryEventTypeCode == CustEventGift || secondaryEventTypeCode == CustEventAmountGift)) ||
		((saleEventTypeCode == SaleEventDiscount || primaryEventTypeCode == CustEventVip) && (saleAmountDtl.SalesmanSaleAmount/saleAmountDtl.TotalListPrice) <= DiscountRate)
}
