package customer

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

type Store struct {
	Id               int64  `json:"id"`
	TenantCode       string `json:"tenantCode"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Manager          string `json:"manager"`
	TelNo            string `json:"telNo"`
	AreaCode         string `json:"areaCode"`
	Area             string `json:"area"`
	Address          string `json:"address"`
	ShippingAreaCode string `json:"shippingAreaCode"`
	ShippingArea     string `json:"shippingArea"`
	ShippingAddress  string `json:"shippingAddress"`
	StatusCode       string `json:"statusCode"`
	Cashier          bool   `json:"cashier"`
	ContractNo       string `json:"contractNo"`
	OpenDate         string `json:"openDate"`
	CloseDate        string `json:"closeDate"`
	Remark           Remark `json:"remark"`
	Enable           bool   `json:"enable"`
}
type Remark struct {
	ElandShopInfos []ElandShopInfo `json:"elandShopInfos,omitempty "`
}
type ElandShopInfo struct {
	BrandId   int64  `json:"brandId,omitempty "`
	BrandCode string `json:"brandCode,omitempty "`
	IsChief   bool   `json:"isChief,omitempty "`
	ShopCode  string `json:"shopCode,omitempty "`
}

func (Store) GetBrandIdByStoreId(ctx context.Context, storeId int64) (int64, error) {
	var resp struct {
		Result struct {
			Items      []Store `json:"items"`
			TotalCount int64   `json:"totalCount"`
		} `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/store/getallinfo?storeIds=%v&withBrand=true&maxResultCount=1",
		config.Config().Services.PlaceManagementApi, storeId)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return 0, err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get store")
		return 0, fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}
	var brandId int64
	if len(resp.Result.Items) > 0 {
		for _, shopInfo := range resp.Result.Items[0].Remark.ElandShopInfos {
			if shopInfo.IsChief {
				brandId = shopInfo.BrandId
			}
		}
	}

	return brandId, nil
}
