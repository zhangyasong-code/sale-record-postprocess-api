package promotion

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"strings"

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
			Code    int64  `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/promotions", config.Config().Services.PromotionApi)
	if _, err := httpreq.New(http.MethodPost, url, p).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("%d-%s-%s", resp.Error.Code, resp.Error.Message, resp.Error.Details)
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
func getStores(ctx context.Context, ids []string) ([]Store, error) {
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
	url := fmt.Sprintf("%s/v1/store/getallinfo?maxResultCount=10&skipCount=0&storeIds=%s&withBrand=true", config.Config().Services.StoreApi, strings.Join(ids, ","))
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
	return resp.Result.Items, nil
}
func reUploadOffer(ctx context.Context, no string) error {
	var resp struct {
		Result interface{} `json:"result"`
		Error  struct {
			Code    int64       `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
		Success bool `json:"success"`
	}
	strArr := strings.Split(no, "-")
	var url string
	if strArr[0] == "1" {
		url = fmt.Sprintf("%s/v1/catalog/campaigns/%s/upload", config.Config().Services.OfferApi, strArr[1])
	} else {
		url = fmt.Sprintf("%s/v1/cart/campaigns/%s/upload", config.Config().Services.OfferApi, strArr[1])
	}
	if _, err := httpreq.New(http.MethodPost, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}
	return nil
}
