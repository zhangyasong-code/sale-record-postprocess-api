package promotion

import (
	"context"
	"errors"
)

type CampaignEventHandler struct {
}

func (h CampaignEventHandler) HandleCartCampaign(ctx context.Context, c CartCampaign) error {
	//1.查询ruleGroup
	ruleGroup, err := getCartRuleGroup(ctx, c.RulesetGroupId)
	if err != nil {
		return err
	}

	//2.转换结构
	promotions, err := CartToCSLEvent(ctx, c, ruleGroup)
	if err != nil {
		return err
	}

	//3.调用promotion-api(上传数据到CSL，并获取eventNo)
	for i := range promotions {
		eventNo, error := getEventNoByPromotion(ctx, promotions[i])
		if err != nil {
			return error
		}
		if eventNo == "" {
			return errors.New("eventNo is null")
		}
		promotions[i].EventNo = eventNo
	}

	//4.保存到数据库
	if err = (PromotionEvent{}).createInArrary(ctx, promotions); err != nil {
		return err
	}

	return nil
}

func (h CampaignEventHandler) HandleCatalogCampaign(ctx context.Context, c CatalogCampaign) error {
	//1.查询ruleset
	ruleset, err := getCatalogRuleGroup(ctx, c.RulesetId)
	if err != nil {
		return err
	}
	//2.转换结构

	promotions, err := CatalogToCSLEvent(ctx, c, ruleset)
	if err != nil {
		return err
	}

	//3.调用promotion-api
	for i := range promotions {
		eventNo, error := getEventNoByPromotion(ctx, promotions[i])
		if err != nil {
			return error
		}
		if eventNo == "" {
			return errors.New("eventNo is null")
		}
		promotions[i].EventNo = eventNo
	}

	//4.保存到数据库(上传数据到CSL，并获取eventNo)
	if err = (PromotionEvent{}).createInArrary(ctx, promotions); err != nil {
		return err
	}

	return nil
}
