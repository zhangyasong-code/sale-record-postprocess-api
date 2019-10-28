package controllers

import (
	"fmt"

	"github.com/labstack/echo"
)

const (
	FlashName      = "flash"
	FlashSeparator = ";"
)

type ApiResult struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   ApiError    `json:"error"`
}

type ApiError struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type ArrayResult struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"totalCount"`
}

var (
	// System Error
	ApiErrorSystem             = ApiError{Code: 10001, Message: "System Error"}
	ApiErrorServiceUnavailable = ApiError{Code: 10002, Message: "Service unavailable"}
	ApiErrorRemoteService      = ApiError{Code: 10003, Message: "Remote service error"}
	ApiErrorIPLimit            = ApiError{Code: 10004, Message: "IP limit"}
	ApiErrorPermissionDenied   = ApiError{Code: 10005, Message: "Permission denied"}
	ApiErrorIllegalRequest     = ApiError{Code: 10006, Message: "Illegal request"}
	ApiErrorHTTPMethod         = ApiError{Code: 10007, Message: "HTTP method is not suported for this request"}
	ApiErrorParameter          = ApiError{Code: 10008, Message: "Parameter error"}
	ApiErrorMissParameter      = ApiError{Code: 10009, Message: "Miss required parameter"}
	ApiErrorDB                 = ApiError{Code: 10010, Message: "DB error, please contact the administator"}
	ApiErrorTokenInvaild       = ApiError{Code: 10011, Message: "Token invaild"}
	ApiErrorMissToken          = ApiError{Code: 10012, Message: "Miss token"}
	ApiErrorVersion            = ApiError{Code: 10013, Message: "API version %s invalid"}
	ApiErrorNotFound           = ApiError{Code: 10014, Message: "Resource not found"}
	// Business Error
	ApiErrorUserNotExists = ApiError{Code: 20001, Message: "User does not exists"}
	ApiErrorPassword      = ApiError{Code: 20002, Message: "Password error"}
)

const defaultMaxResultCount = 20

func ReturnApiFail(c echo.Context, status int, apiError ApiError, err error, v ...interface{}) error {
	str := ""
	if err != nil {
		str = err.Error()
	}
	return c.JSON(status, ApiResult{
		Success: false,
		Error: ApiError{
			Code:    apiError.Code,
			Message: fmt.Sprintf(apiError.Message, v...),
			Details: str,
		},
	})
}

func ReturnApiSucc(c echo.Context, status int, result interface{}) error {

	return c.JSON(status, ApiResult{
		Success: true,
		Result:  result,
	})
}
