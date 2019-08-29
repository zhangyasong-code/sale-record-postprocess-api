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
}

func (PromotionEventController) GetOne(c echo.Context) error {
	no := c.Param("no")

	if no == "" {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, errors.New("no can not be null"))
	}

	v, err := promotion.GetById(c.Request().Context(), no)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if v == nil {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}
