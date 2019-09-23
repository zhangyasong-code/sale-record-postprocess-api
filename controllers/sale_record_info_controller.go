package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"nhub/sale-record-postprocess-api/customer"
	"nhub/sale-record-postprocess-api/factory"
	"nhub/sale-record-postprocess-api/payamt"
	"nhub/sale-record-postprocess-api/promotion"
	"nhub/sale-record-postprocess-api/salePerson"
	"nhub/sale-record-postprocess-api/saleRecordFee"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type SaleRecordInfoController struct{}

func (c SaleRecordInfoController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("/:transactionId", c.GetOne).
		AddParamPath("", "transactionId", "transactionId")
}

func (SaleRecordInfoController) GetOne(c echo.Context) error {
	transactionId, _ := strconv.ParseInt(c.Param("transactionId"), 10, 64)
	if transactionId == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, errors.New("invalid transactionId"))
	}

	v, err := getSaleRecordInfo(c.Request().Context(), transactionId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}

func getSaleRecordInfo(ctx context.Context, transactionId int64) (interface{}, error) {
	var postMileages []customer.PostMileage
	if err := factory.SaleRecordDB(ctx).Where("transaction_id = ?", transactionId).Find(&postMileages); err != nil {
		return nil, err
	}

	postMileageIds := []int64{}
	for _, postMileage := range postMileages {
		postMileageIds = append(postMileageIds, postMileage.Id)
	}
	var postMileageDtls []customer.PostMileageDtl
	if len(postMileageIds) != 0 {
		if err := factory.SaleRecordDB(ctx).Where("1=1").In("post_mileage_id", postMileageDtls).Find(&postMileageDtls); err != nil {
			return nil, err
		}
	}

	var saleRecordDtlSalesmanAmounts []salePerson.SaleRecordDtlSalesmanAmount
	if err := factory.SaleRecordDB(ctx).Where("transaction_id = ?", transactionId).Find(&saleRecordDtlSalesmanAmounts); err != nil {
		return nil, err
	}

	var itemOffers, cartOffers []struct {
		CouponNo                 string  `json:"couponNo,omitempty"`
		ItemCode                 string  `json:"itemCode,omitempty"`
		ItemCodes                string  `json:"itemCodes,omitempty"`
		Price                    float64 `json:"price,omitempty"`
		TenantCode               string  `json:"tenantCode,omitempty"`
		Type                     string  `json:"type,omitempty"`
		promotion.PromotionEvent `xorm:"extends"`
	}

	if err := factory.SaleRecordDB(ctx).
		Table("applied_sale_record_item_offer").Alias("o").Select("o.*, p.*").
		Join("INNER", []string{"assorted_sale_record_dtl", "d"}, "o.transaction_dtl_id = d.id").
		Join("LEFT", []string{"promotion_event", "p"}, "p.offer_no = o.offer_no").
		Where("d.transaction_id = ?", transactionId).
		Find(&itemOffers); err != nil {
		return nil, err
	}

	if err := factory.SaleRecordDB(ctx).Table("applied_sale_record_cart_offer").Alias("o").Select("o.*, p.*").
		Join("LEFT", []string{"promotion_event", "p"}, "p.offer_no = o.offer_no").
		Where("o.transaction_id = ?", transactionId).
		Find(&cartOffers); err != nil {
		return nil, err
	}

	var postSaleRecordFees []saleRecordFee.PostSaleRecordFee
	if err := factory.SaleRecordDB(ctx).Where("transaction_id = ?", transactionId).Find(&postSaleRecordFees); err != nil {
		return nil, err
	}

	var postFailCreateSaleFees []saleRecordFee.PostFailCreateSaleFee
	if err := factory.SaleRecordDB(ctx).Where("transaction_id = ?", transactionId).Find(&postFailCreateSaleFees); err != nil {
		return nil, err
	}

	var postPayments []payamt.PostPayment
	if err := factory.SaleRecordDB(ctx).Where("transaction_id = ?", transactionId).Find(&postPayments); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"transactionId":                transactionId,
		"postMileage":                  postMileages,
		"postMileageDtls":              postMileageDtls,
		"saleRecordDtlSalesmanAmounts": saleRecordDtlSalesmanAmounts,
		"itemOffers":                   itemOffers,
		"cartOffers":                   cartOffers,
		"postSaleRecordFees":           postSaleRecordFees,
		"postFailCreateSaleFees":       postFailCreateSaleFees,
		"postPayments":                 postPayments,
	}, nil
}
