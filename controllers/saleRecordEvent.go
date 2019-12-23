package controllers

import (
	"encoding/json"
	"nhub/sale-record-postprocess-api/customer"
	"nhub/sale-record-postprocess-api/models"
	"nhub/sale-record-postprocess-api/postprocess"
	"nhub/sale-record-postprocess-api/salePerson"
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
	recordStr, _ := json.Marshal(event)

	if err := (customer.CustomerEventHandler{}).Handle(ctx, event); err != nil {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			OrderId:      event.OrderId,
			RefundId:     event.RefundId,
			StoreId:      event.StoreId,
			ModuleType:   string(postprocess.ModuleMileage),
			IsSuccess:    false,
			Error:        err.Error(),
			ModuleEntity: string(recordStr),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, saveErr)
		}
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	} else {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			OrderId:      event.OrderId,
			RefundId:     event.RefundId,
			StoreId:      event.StoreId,
			ModuleType:   string(postprocess.ModuleMileage),
			IsSuccess:    true,
			Error:        "",
			ModuleEntity: string(recordStr),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, saveErr)
		}
	}

	if err := (salePerson.SalesPersonEventHandler{}).Handle(ctx, event); err != nil {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			OrderId:      event.OrderId,
			RefundId:     event.RefundId,
			StoreId:      event.StoreId,
			ModuleType:   string(postprocess.ModuleSalePerson),
			IsSuccess:    false,
			Error:        err.Error(),
			ModuleEntity: string(recordStr),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, saveErr)
		}
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	} else {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			OrderId:      event.OrderId,
			RefundId:     event.RefundId,
			StoreId:      event.StoreId,
			ModuleType:   string(postprocess.ModuleSalePerson),
			IsSuccess:    true,
			Error:        "",
			ModuleEntity: string(recordStr),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, saveErr)
		}
	}

	if err := (saleRecordFee.SaleRecordFeeEventHandler{}).Handle(ctx, event); err != nil {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			OrderId:      event.OrderId,
			RefundId:     event.RefundId,
			StoreId:      event.StoreId,
			ModuleType:   string(postprocess.ModuleSaleFee),
			IsSuccess:    false,
			Error:        err.Error(),
			ModuleEntity: string(recordStr),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, saveErr)
		}
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	} else {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			OrderId:      event.OrderId,
			RefundId:     event.RefundId,
			StoreId:      event.StoreId,
			ModuleType:   string(postprocess.ModuleSaleFee),
			IsSuccess:    true,
			Error:        "",
			ModuleEntity: string(recordStr),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, saveErr)
		}
	}

	return ReturnApiSucc(c, http.StatusOK, event)
}
