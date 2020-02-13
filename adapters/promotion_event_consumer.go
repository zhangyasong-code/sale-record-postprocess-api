package adapters

import (
	"encoding/json"
	"nhub/sale-record-postprocess-api/promotion"
	"nomni/utils/eventconsume"

	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
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

	str, _ := json.Marshal(event)
	logrus.WithField("Body", string(str)).Info("Offer Event Body>>>>>>")

	if c.Status() == EventCartCampaignApproved {
		var event promotion.CartCampaign

		if err := c.Bind(&event); err != nil {
			return err
		}

		if err := (promotion.CampaignEventHandler{}).HandleCartCampaign(ctx, event); err != nil {
			logrus.WithFields(logrus.Fields{
				"Id":              event.Id,
				"RulesetGroupId ": event.RulesetGroupId,
				"Error":           err,
			}).Info("Fail to handle event")
			return err
		}

		logrus.WithFields(logrus.Fields{
			"Id":              event.Id,
			"RulesetGroupId ": event.RulesetGroupId,
		}).Info("Success to handle event")
	}

	if c.Status() == EventCatalogCampaignApproved {
		var event promotion.CatalogCampaign

		if err := c.Bind(&event); err != nil {
			return err
		}

		if err := (promotion.CampaignEventHandler{}).HandleCatalogCampaign(ctx, event); err != nil {
			logrus.WithFields(logrus.Fields{
				"Id":              event.Id,
				"RulesetGroupId ": event.RulesetId,
				"Error":           err,
			}).Info("Success to handle event")
			return err
		}

		logrus.WithFields(logrus.Fields{
			"Id":              event.Id,
			"RulesetGroupId ": event.RulesetId,
		}).Info("Success to handle event")
	}

	return nil
}
