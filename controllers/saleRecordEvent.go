package controllers

import (
	"nhub/sale-record-postprocess-api/models"
	"nhub/sale-record-postprocess-api/saleRecordFee"

	"github.com/pangpanglabs/echoswagger"

	"net/http"

	"github.com/labstack/echo"
)

type SaleRecordEventController struct{}

func (c SaleRecordEventController) Init(g echoswagger.ApiGroup) {
	g.POST("", c.HandleEvent)
}

func (SaleRecordEventController) HandleEvent(c echo.Context) error {
	var event models.SaleRecordEvent
	if err := c.Bind(&event); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	ctx := c.Request().Context()
	if err := (models.CustomerEventHandler{}).Handle(ctx, event); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	if err := (models.SalesPersonEventHandler{}).Handle(ctx, event); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	if err := (saleRecordFee.SaleRecordFeeEventHandler{}).Handle(ctx, event); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	return ReturnApiSucc(c, http.StatusOK, event)
}
