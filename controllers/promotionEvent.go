package controllers

import (
	"net/http"
	"nhub/sale-record-postprocess-api/models"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type PromotionEventController struct{}

func (c PromotionEventController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("/:id", c.GetOne).
		AddParamPath("", "id", "Id of Offer")
}

func (PromotionEventController) GetOne(c echo.Context) error {
	arr := strings.Split(c.Param("id"), "-")
	campaignId, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	ruleId, err := strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	v, err := models.GetById(c.Request().Context(), campaignId, ruleId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if v == nil {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	return ReturnApiSucc(c, http.StatusOK, v)
	return nil
}
