package salePerson

import (
	"context"
	"math"
	"nhub/sale-record-postprocess-api/customer"
	"nhub/sale-record-postprocess-api/models"
	"nhub/sale-record-postprocess-api/promotion"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"github.com/sirupsen/logrus"
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

func (h SalesPersonEventHandler) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.POST("", h.HandleTest)
}

func (h SalesPersonEventHandler) Handle(ctx context.Context, s models.SaleRecordEvent) error {

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
			TotalDiscountItemOfferPrice: s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedItemOfferPrice,
			TotalDiscountCartOfferPrice: s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedCartOfferPrice,
			TotalPaymentPrice:           s.AssortedSaleRecordDtlList[i].TotalPrice.ListPrice - s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedCartOfferPrice - s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedItemOfferPrice,
			TransactionType:             s.TransactionType,
			TransactionChannelType:      s.TransactionChannelType,
			SalesmanSaleDiscountRate:    0,
			SalesmanSaleAmount:          s.AssortedSaleRecordDtlList[i].TotalPrice.ListPrice - s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedCartOfferPrice - s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedItemOfferPrice,
			TransactionCreateDate:       s.TransactionCreateDate,
		}
		//查询使用积分
		// var channelType customer.UseType
		// if strings.ToUpper(s.TransactionType) == "MINUS" {
		// 	channelType = customer.UseTypeUsed
		// } else {
		// 	channelType = customer.UseTypeUsedCancel
		// }
		// mileage, mileagePrice, err := GetUsedBonus(ctx, s.AssortedSaleRecordDtlList[i].Id, s.AssortedSaleRecordDtlList[i].OrderItemId, s.AssortedSaleRecordDtlList[i].RefundItemId, channelType)
		// if err != nil {
		// 	logrus.WithField("err", err).Info("GetPostMileageDtlError")
		// 	return err
		// }
		saleAmountDtl.Mileage = s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedPaymentPrice - s.AssortedSaleRecordDtlList[i].DistributedPrice.DistributedCashPrice
		saleAmountDtl.MileagePrice = s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedPaymentPrice - s.AssortedSaleRecordDtlList[i].DistributedPrice.DistributedCashPrice
		//计算营业员业绩金额-SalesmanSaleAmount
		offers := []SaleRecordDtlOffer{}
		if s.TotalPrice.DiscountPrice == 0 {
			saleAmountDtl.SalesmanSaleAmount = saleAmountDtl.TotalPaymentPrice
			saleAmountDtl.SalesmanNormalSaleAmount = saleAmountDtl.TotalPaymentPrice
			saleAmountDtl.SalesmanDiscountSaleAmount = 0
		} else {
			//查询使用优惠的类型
			offers, salesmanSaleDiscountRate, salesmanSaleAmount, err := GetDiscountType(ctx, s.AssortedSaleRecordDtlList[i].ItemOffers, s.TransactionType, saleAmountDtl.TotalPaymentPrice, math.Abs(saleAmountDtl.MileagePrice))
			if err != nil {
				logrus.WithField("err", err).Info("GetPostMileageDtlError")
				return err
			}
			saleAmountDtl.SalesmanSaleDiscountRate = salesmanSaleDiscountRate
			if salesmanSaleAmount != 0 {
				saleAmountDtl.SalesmanSaleAmount = salesmanSaleAmount
			}

			//拆分业绩金额-正常业绩和折扣业绩
			normalSaleAmount, discountSaleAmount := SeparateNormalAndDiscountAmt(offers, saleAmountDtl.SalesmanSaleAmount, saleAmountDtl.TotalListPrice)
			saleAmountDtl.SalesmanNormalSaleAmount = normalSaleAmount
			saleAmountDtl.SalesmanDiscountSaleAmount = discountSaleAmount
		}

		//保存数据
		dtl, err := SaveAchievement(ctx, saleAmountDtl)
		if err != nil {
			logrus.WithField("err", err).Info("SaveSaleAmountDtlFail")
			return err
		}

		if err := SaveOffer(ctx, offers, dtl.Id); err != nil {
			logrus.WithField("err", err).Info("SaveSaleAmountOfferFail")
			return err
		}
	}
	return nil
}

func isNormalAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode string, salesmanSaleAmount, totalListPrice float64) bool {
	return (primaryEventTypeCode == CustEventBrandDiscount || primaryEventTypeCode == CustEventCoupon || primaryEventTypeCode == CustEventCustmoer) ||
		((saleEventTypeCode == SaleEventDiscount || primaryEventTypeCode == CustEventVip) && (salesmanSaleAmount/totalListPrice) > DiscountRate) ||
		(saleEventTypeCode == "" && primaryEventTypeCode == "" && (secondaryEventTypeCode == CustEventMileage || secondaryEventTypeCode == ""))
}
func isDiscountAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode string, salesmanSaleAmount, totalListPrice float64) bool {
	return ((saleEventTypeCode == SaleEventGive || saleEventTypeCode == SaleEventReduce) || (secondaryEventTypeCode == CustEventGift || secondaryEventTypeCode == CustEventAmountGift)) ||
		((saleEventTypeCode == SaleEventDiscount || primaryEventTypeCode == CustEventVip) && (salesmanSaleAmount/totalListPrice) <= DiscountRate)
}
func GetUsedBonus(ctx context.Context, dtlId, dtlOrderItemId, dtlRefundItemId int64, useType customer.UseType) (float64, float64, error) {
	var mileage float64
	var mileagePrice float64
	has, e, err := customer.PostMileageDtl{}.GetByKey(ctx, dtlId, dtlOrderItemId, dtlRefundItemId, useType)
	if err != nil {
		return 0, 0, err
	}
	if has {
		mileage = e.Point
		mileagePrice = e.PointPrice
	} else {
		mileage = 0
		mileagePrice = 0
	}
	return mileage, mileagePrice, nil
}
func SeparateNormalAndDiscountAmt(offers []SaleRecordDtlOffer, salesmanSaleAmount, totalListPrice float64) (float64, float64) {
	var normalSaleAmount float64
	var discountSaleAmount float64
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
	if isNormalAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode, salesmanSaleAmount, totalListPrice) {
		normalSaleAmount = salesmanSaleAmount
		discountSaleAmount = 0
	} else if isDiscountAmt(primaryEventTypeCode, secondaryEventTypeCode, saleEventTypeCode, salesmanSaleAmount, totalListPrice) {
		normalSaleAmount = 0
		discountSaleAmount = salesmanSaleAmount
	}
	return normalSaleAmount, discountSaleAmount
}
func GetDiscountType(ctx context.Context, offers []models.Offer, channelType string, totalPaymentPrice, mileagePrice float64) ([]SaleRecordDtlOffer, float64, float64, error) {
	res := []SaleRecordDtlOffer{}
	var salesmanSaleDiscountRate, salesmanSaleAmount float64
	for n := 0; n < len(offers); n++ {
		p, err := promotion.GetByNo(ctx, offers[n].OfferNo)
		if err != nil {
			return res, 0, 0, err
		}
		offer := SaleRecordDtlOffer{
			OfferId:         offers[n].OfferId,
			EventType:       p.EventTypeCode,
			SaleBaseAmt:     p.SaleBaseAmt,
			DiscountBaseAmt: p.DiscountBaseAmt,
			DiscountRate:    p.DiscountRate,
		}
		res = append(res, offer)
		//计算实际销售额活动扣率
		if p.EventTypeCode == SaleEventGive { //送活动
			salesmanSaleDiscountRate = math.Floor((p.DiscountBaseAmt/(p.SaleBaseAmt*0.1+p.SaleBaseAmt+p.DiscountBaseAmt))*100) / 100
		}
		switch p.EventTypeCode {
		case SaleEventReduce, SaleEventDiscount: //减 折
			salesmanSaleAmount = totalPaymentPrice - mileagePrice
		case SaleEventGive: //送
			salesmanSaleAmount = math.Floor((totalPaymentPrice-mileagePrice-totalPaymentPrice*salesmanSaleDiscountRate)*100) / 100
		}
	}
	return res, salesmanSaleDiscountRate, salesmanSaleAmount, nil
}
func SaveAchievement(ctx context.Context, dtl SaleRecordDtlSalesmanAmount) (SaleRecordDtlSalesmanAmount, error) {
	has, _, err := (&dtl).GetByKey(ctx, dtl.TransactionId, dtl.SaleRecordDtlId)
	if err != nil {
		logrus.WithField("err", err).Info("CheckSaleAmountDtlExist")
		return dtl, err
	}
	if has {
		_, err := (&dtl).Update(ctx, dtl.TransactionId, dtl.SaleRecordDtlId)
		if err != nil {
			logrus.WithField("err", err).Info("UpdateSaleAmountDtlError")
			return dtl, err
		}
	} else {
		if err := (&dtl).Create(ctx); err != nil {
			logrus.WithField("err", err).Info("CreateSaleAmountDtlError")
			return dtl, err
		}
	}
	return dtl, nil
}
func SaveOffer(ctx context.Context, offers []SaleRecordDtlOffer, saleAmountDtlId int64) error {
	for e := 0; e < len(offers); e++ {
		offers[e].SalesmanAmountId = saleAmountDtlId
		has, _, err := (&offers[e]).GetByKey(ctx, offers[e].OfferId, offers[e].SalesmanAmountId)
		if err != nil {
			logrus.WithField("err", err).Info("CheckSaleAmountOfferExist")
			return err
		}
		if has {
			_, err := (&offers[e]).Update(ctx, offers[e].OfferId, offers[e].SalesmanAmountId)
			if err != nil {
				logrus.WithField("err", err).Info("UpdateSaleAmountOfferError")
				return err
			}
		} else {
			if err := (&offers[e]).Create(ctx); err != nil {
				logrus.WithField("err", err).Info("CreateSaleAmountOfferError")
				return err
			}
		}
	}
	return nil
}
func (h SalesPersonEventHandler) HandleTest(c echo.Context) error {
	var s models.SaleRecordEvent
	if err := c.Bind(&s); err != nil {
		return err
	}

	ctx := c.Request().Context()

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
			TransactionType:             s.TransactionType,
			TransactionChannelType:      s.TransactionChannelType,
			SalesmanSaleDiscountRate:    0,
			SalesmanSaleAmount:          s.AssortedSaleRecordDtlList[i].DistributedPrice.TotalDistributedPaymentPrice,
			TransactionCreateDate:       s.TransactionCreateDate,
		}
		//查询使用积分
		var channelType customer.UseType
		if s.TransactionType == "MINUS" {
			channelType = customer.UseTypeUsed
		} else {
			channelType = customer.UseTypeUsedCancel
		}
		mileage, mileagePrice, err := GetUsedBonus(ctx, s.AssortedSaleRecordDtlList[i].Id, s.AssortedSaleRecordDtlList[i].OrderItemId, s.AssortedSaleRecordDtlList[i].RefundItemId, channelType)
		if err != nil {
			logrus.WithField("err", err).Info("GetPostMileageDtlError")
			return err
		}
		saleAmountDtl.Mileage = mileage
		saleAmountDtl.MileagePrice = mileagePrice
		//计算营业员业绩金额-SalesmanSaleAmount
		offers := []SaleRecordDtlOffer{}
		if s.TotalPrice.DiscountPrice == 0 {
			saleAmountDtl.SalesmanSaleAmount = saleAmountDtl.TotalPaymentPrice
			saleAmountDtl.SalesmanNormalSaleAmount = saleAmountDtl.TotalPaymentPrice
			saleAmountDtl.SalesmanDiscountSaleAmount = 0
		} else {
			//查询使用优惠的类型
			offers, salesmanSaleDiscountRate, salesmanSaleAmount, err := GetDiscountType(ctx, s.AssortedSaleRecordDtlList[i].ItemOffers, s.TransactionType, saleAmountDtl.TotalPaymentPrice, saleAmountDtl.MileagePrice)
			if err != nil {
				logrus.WithField("err", err).Info("GetPostMileageDtlError")
				return err
			}
			saleAmountDtl.SalesmanSaleDiscountRate = salesmanSaleDiscountRate
			saleAmountDtl.SalesmanSaleAmount = salesmanSaleAmount

			//拆分业绩金额-正常业绩和折扣业绩
			normalSaleAmount, discountSaleAmount := SeparateNormalAndDiscountAmt(offers, saleAmountDtl.SalesmanSaleAmount, saleAmountDtl.TotalListPrice)
			saleAmountDtl.SalesmanNormalSaleAmount = normalSaleAmount
			saleAmountDtl.SalesmanDiscountSaleAmount = discountSaleAmount
		}

		//保存数据
		dtl, err := SaveAchievement(ctx, saleAmountDtl)
		if err != nil {
			logrus.WithField("err", err).Info("SaveSaleAmountDtlFail")
			return err
		}

		if err := SaveOffer(ctx, offers, dtl.Id); err != nil {
			logrus.WithField("err", err).Info("SaveSaleAmountOfferFail")
			return err
		}
	}
	return nil
}
