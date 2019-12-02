package controllers

import (
	"net/http"
	"nhub/sale-record-postprocess-api/postprocess"
	"nomni/utils/api"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type ProcessLogController struct{}

func (c ProcessLogController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("/PostProcessFail", c.GetPostProcessFail)
	g.GET("/PostProcessFails", c.GetPostProcessFails)
	g.PUT("/:id", c.UpdatePostProcessSuccess)
}

func (ProcessLogController) GetPostProcessFail(c echo.Context) error {
	transactionId, _ := strconv.ParseInt(c.QueryParam("transactionId"), 10, 64)
	if transactionId == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.MissRequiredParamError("transactionId"))
	}
	moduleType := c.QueryParam("moduleType")
	if moduleType == "" {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.MissRequiredParamError("moduleType"))
	}

	v, err := postprocess.PostProcessSuccess{}.Get(c.Request().Context(), false, transactionId, moduleType)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}

func (ProcessLogController) GetPostProcessFails(c echo.Context) error {
	transactionId, _ := strconv.ParseInt(c.QueryParam("transactionId"), 10, 64)
	orderId, _ := strconv.ParseInt(c.QueryParam("orderId"), 10, 64)
	refundId, _ := strconv.ParseInt(c.QueryParam("refundId"), 10, 64)
	moduleType := c.QueryParam("moduleType")
	skipCount, _ := strconv.ParseInt(c.QueryParam("skipCount"), 10, 64)
	maxResultCount, _ := strconv.ParseInt(c.QueryParam("maxResultCount"), 10, 64)

	if maxResultCount == 0 || maxResultCount > 200 {
		maxResultCount = defaultMaxResultCount
	}

	totalCount, result, err := postprocess.PostProcessSuccess{}.GetAll(c.Request().Context(), false, transactionId, orderId, refundId, moduleType, int(skipCount), int(maxResultCount))
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, ArrayResult{
		TotalCount: totalCount,
		Items:      result,
	})
}

func (ProcessLogController) UpdatePostProcessSuccess(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if id == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.MissRequiredParamError("id"))
	}

	var postProcessSuccess postprocess.PostProcessSuccess
	if err := c.Bind(&postProcessSuccess); err != nil {
		ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	postprocess, err := postprocess.PostProcessSuccess{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if postprocess.Id == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.NotFoundError())
	}
	if postprocess.IsSuccess == true {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.NotFoundError())
	}
	if postProcessSuccess.OrderId != postprocess.OrderId || postProcessSuccess.RefundId != postprocess.RefundId || postProcessSuccess.ModuleType != postprocess.ModuleType || postProcessSuccess.TransactionId != postprocess.TransactionId {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.NotFoundError())
	}
	if err := postProcessSuccess.Update(c.Request().Context()); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, nil)
}
