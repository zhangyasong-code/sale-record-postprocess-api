package controllers

import (
	"errors"
	"net/http"
	"nhub/sale-record-postprocess-api/promotion"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type PromotionEventController struct{}

func (c PromotionEventController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("/:no", c.GetOne).
		AddParamPath("", "no", "no of Offer")
	g.GET("", c.GetAll).
		AddParamBody(promotion.SearchInput{}, "", "", true)
	g.GET("/coupon", c.GetCouponEvent).
		AddParamBody(promotion.SearchInput{}, "", "", true)
}

func (PromotionEventController) GetOne(c echo.Context) error {
	no := c.Param("no")

	if no == "" {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, errors.New("no can not be null"))
	}

	v, err := promotion.GetByNo(c.Request().Context(), no)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if v == nil {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}

func (PromotionEventController) GetAll(c echo.Context) error {
	var v promotion.SearchInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	if v.MaxResultCount == 0 {
		v.MaxResultCount = 10
	}
	result, totalCount, err := promotion.PromotionEvent{}.GetAll(c.Request().Context(), v)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, ArrayResult{
		TotalCount: totalCount,
		Items:      result,
	})
}

func (PromotionEventController) GetCouponEvent(c echo.Context) error {
	var v promotion.SearchInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	if v.MaxResultCount == 0 {
		v.MaxResultCount = 10
	}
	result, totalCount, err := promotion.PostCouponEvent{}.GetAll(c.Request().Context(), v)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, ArrayResult{
		TotalCount: totalCount,
		Items:      result,
	})
}
