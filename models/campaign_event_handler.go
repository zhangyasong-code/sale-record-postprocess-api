package models

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
	p, err := CartToCSLEvent(ctx, c, ruleGroup)
	if err != nil {
		return err
	}

	//3.调用promotion-api(上传数据到CSL，并获取eventNo)
	eventNo, error := getEventNoByPromotion(ctx, *p)
	if err != nil {
		return error
	}
	if eventNo == "" {
		return errors.New("eventNo is null")
	}

	p.EventNo = eventNo

	//4.保存到数据库
	if err := p.create(ctx); err != nil {
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

	p, err := CatalogToCSLEvent(ctx, c, ruleset)
	if err != nil {
		return err
	}

	//3.调用promotion-api
	eventNo, error := getEventNoByPromotion(ctx, *p)
	if err != nil {
		return error
	}
	if eventNo == "" {
		return errors.New("eventNo is null")
	}

	p.EventNo = eventNo

	//4.保存到数据库(上传数据到CSL，并获取eventNo)
	if err := p.create(ctx); err != nil {
		return err
	}

	return nil
}
