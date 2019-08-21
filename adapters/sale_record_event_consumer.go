package adapters

import (
	"nhub/sale-record-postprocess-api/models"
	"nomni/utils/eventconsume"

	"github.com/pangpanglabs/goutils/kafka"
)

const (
	EventCatalogCampaignApproved = "CatalogCampaignApproved"
	EventCartCampaignApproved    = "CartCampaignApproved"
)

func NewSaleRecordEventConsumer(serviceName string, kafkaConfig kafka.Config, filters ...eventconsume.Filter) error {
	return eventconsume.NewEventConsumer(serviceName, kafkaConfig.Brokers, kafkaConfig.Topic, filters).Handle(handleEvent)
}

func handleEvent(c eventconsume.ConsumeContext) error {
	ctx := c.Context()

	if c.Status() == EventCartCampaignApproved {
		var event models.CartCampaign

		if err := c.Bind(&event); err != nil {
			return err
		}

		if err := (models.CampaignEventHandler{}).HandleCartCampaign(ctx, event); err != nil {
			return err
		}
	}

	if c.Status() == EventCatalogCampaignApproved {
		var event models.CatalogCampaign

		if err := c.Bind(&event); err != nil {
			return err
		}

		if err := (models.CampaignEventHandler{}).HandleCatalogCampaign(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
