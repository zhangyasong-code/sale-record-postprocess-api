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
}

func (ProcessLogController) GetPostProcessFail(c echo.Context) error {
	orderId, err := strconv.ParseInt(c.QueryParam("orderId"), 10, 64)
	if orderId == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.MissRequiredParamError("orderId"))
	}
	refundId, _ := strconv.ParseInt(c.QueryParam("refundId"), 10, 64)
	moduleType := c.QueryParam("moduleType")
	if moduleType == "" {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.MissRequiredParamError("moduleType"))
	}

	v, err := postprocess.PostProcessSuccess{}.Get(c.Request().Context(), false, orderId, refundId, moduleType)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, v)
}

func (ProcessLogController) GetPostProcessFails(c echo.Context) error {
	orderId, err := strconv.ParseInt(c.QueryParam("orderId"), 10, 64)
	if orderId == 0 {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, api.MissRequiredParamError("orderId"))
	}
	refundId, _ := strconv.ParseInt(c.QueryParam("refundId"), 10, 64)
	moduleType := c.QueryParam("moduleType")
	skipCount, _ := strconv.ParseInt(c.QueryParam("skipCount"), 10, 64)
	maxResultCount, _ := strconv.ParseInt(c.QueryParam("maxResultCount"), 10, 64)

	if maxResultCount == 0 || maxResultCount > 200 {
		maxResultCount = defaultMaxResultCount
	}

	totalCount, result, err := postprocess.PostProcessSuccess{}.GetAll(c.Request().Context(), false, orderId, refundId, moduleType, int(skipCount), int(maxResultCount))
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	return ReturnApiSucc(c, http.StatusOK, ArrayResult{
		TotalCount: totalCount,
		Items:      result,
	})
}
