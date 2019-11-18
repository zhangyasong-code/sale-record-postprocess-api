package customer

import (
	"context"
	"fmt"
	"net/http"
	"nhub/sale-record-postprocess-api/config"
	"strconv"
	"strings"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/sirupsen/logrus"
)

type BrandInfo struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Product struct {
	Id         int64             `json:"id"`
	Code       string            `json:"code"`
	Name       string            `json:"name"`
	Brand      BrandInfo         `json:"brand"`
	TitleImage string            `json:"titleImage"`
	ListPrice  float64           `json:"listPrice"`
	Attributes map[string]string `json:"attributes,omitempty"`
	HasDigital bool              `json:"hasDigital"`
	Enable     bool              `json:"enable"`
}

func (BrandInfo) GetBrandInfo(ctx context.Context, brandId int64) (string, error) {
	var resp struct {
		Result  *BrandInfo `json:"result"`
		Success bool       `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}

	url := fmt.Sprintf("%s/v1/brands/%v", config.Config().Services.ProductApi, brandId)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return "", err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get brandcode")
		return "", fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}
	return resp.Result.Code, nil
}

func (BrandInfo) GetCustBrandId(ctx context.Context, brandIds string) (int64, error) {
	/*获取品牌Id*/

	brandMembers, err := BrandMember{}.GetBrandMembers(ctx, brandIds)
	if err != nil {
		return 0, err
	}
	memberBrandIdsStr := ""
	for _, brandMember := range brandMembers {
		tempStr := "," + strconv.FormatInt(brandMember.MemberBrandId, 10)
		if !strings.Contains(memberBrandIdsStr, tempStr) {
			memberBrandIdsStr = memberBrandIdsStr + tempStr
		}
	}
	memberBrandIds := strings.Split(memberBrandIdsStr, ",")
	if len(memberBrandIds) != 2 {
		return 0, fmt.Errorf("Members of multiple brands are not supported")
	}
	memberBrandId, _ := strconv.ParseInt(memberBrandIds[1], 10, 64)
	return memberBrandId, nil
}

func (Product) GetBrandIdsByProductIds(ctx context.Context, productIds string) (string, error) {
	var resp struct {
		Result struct {
			Items      []Product `json:"items"`
			TotalCount int64     `json:"totalCount"`
		} `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		} `json:"error"`
	}
	url := fmt.Sprintf("%s/v1/products?ids=%s", config.Config().Services.ProductApi, productIds)
	logrus.WithField("url", url).Info("url")
	if _, err := httpreq.New(http.MethodGet, url, nil).
		WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
		Call(&resp); err != nil {
		return "", err
	} else if !resp.Success {
		logrus.WithFields(logrus.Fields{
			"errorCode":    resp.Error.Code,
			"errorMessage": resp.Error.Message,
		}).Error("Fail to get brandmembers")
		return "", fmt.Errorf("[%d]%s", resp.Error.Code, resp.Error.Details)
	}

	brandIds := ""
	for _, p := range resp.Result.Items {
		brandIds = brandIds + strconv.FormatInt(p.Brand.Id, 10) + ","
	}

	return brandIds, nil
}
