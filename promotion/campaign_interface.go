package promotion

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
)

func getCartRuleGroup(ctx context.Context, id int64) (*CartRulesetGroup, error) {
	var resp struct {
		Result  CartRulesetGroup `json:"result"`
		Success bool             `json:"success"`
		Error   struct {
			Code    int64       `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/cart/ruleset-groups/%d", config.Config().Services.OfferApi, id)
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}
	return &resp.Result, nil
}

func getCatalogRuleGroup(ctx context.Context, id int64) (*CatalogRuleset, error) {
	var resp struct {
		Result  CatalogRuleset `json:"result"`
		Success bool           `json:"success"`
		Error   struct {
			Code    int64       `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/catalog/rulesets/%d", config.Config().Services.OfferApi, id)
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}
	return &resp.Result, nil
}

func getEventNoByPromotion(ctx context.Context, p PromotionEvent) (string, error) {
	var resp struct {
		Result  string `json:"result"`
		Success bool   `json:"success"`
		Error   struct {
			Code    int64       `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/promotions", config.Config().Services.PromotionApi)
	if _, err := httpreq.New(http.MethodPost, url, p).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}
	return resp.Result, nil
}

func getBrand(ctx context.Context, id int64) (*Brand, error) {
	var resp struct {
		Result  Brand `json:"result"`
		Success bool  `json:"success"`
		Error   struct {
			Code    int64       `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/brands/%d", config.Config().Services.ProductApi, id)
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}
	return &resp.Result, nil
}
func getStore(ctx context.Context, id int64) (*Store, error) {
	var resp struct {
		Result struct {
			TotalCount int     `json:"totalCount"`
			Items      []Store `json:"items"`
		} `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int64       `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/store/getallinfo?maxResultCount=10&skipCount=0&storeIds=%d", config.Config().Services.StoreApi, id)
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}
	if len(resp.Result.Items) == 0 {
		return nil, errors.New("Store is not exist")
	}
	return &resp.Result.Items[0], nil
}