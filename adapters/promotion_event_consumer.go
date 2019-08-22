package adapters

import (
	"nhub/sale-record-postprocess-api/promotion"
	"nomni/utils/eventconsume"

	"github.com/pangpanglabs/goutils/kafka"
)

const (
	EventCatalogCampaignApproved = "CatalogCampaignApproved"
	EventCartCampaignApproved    = "CartCampaignApproved"
)

func NewPromotionEventConsumer(serviceName string, kafkaConfig kafka.Config, filters ...eventconsume.Filter) error {
	return eventconsume.NewEventConsumer(serviceName, kafkaConfig.Brokers, kafkaConfig.Topic, filters).Handle(handlePromotionEvent)
}

func handlePromotionEvent(c eventconsume.ConsumeContext) error {
	ctx := c.Context()

	if c.Status() == EventCartCampaignApproved {
		var event promotion.CartCampaign

		if err := c.Bind(&event); err != nil {
			return err
		}

		if err := (promotion.CampaignEventHandler{}).HandleCartCampaign(ctx, event); err != nil {
			return err
		}
	}

	if c.Status() == EventCatalogCampaignApproved {
		var event promotion.CatalogCampaign

		if err := c.Bind(&event); err != nil {
			return err
		}

		if err := (promotion.CampaignEventHandler{}).HandleCatalogCampaign(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
