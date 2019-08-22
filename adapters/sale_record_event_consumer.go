package adapters

import (
	"nhub/sale-record-postprocess-api/models"
	"nomni/utils/eventconsume"

	"github.com/pangpanglabs/goutils/kafka"
)

func NewSaleRecordEventConsumer(serviceName string, kafkaConfig kafka.Config, filters ...eventconsume.Filter) error {
	return eventconsume.NewEventConsumer(serviceName, kafkaConfig.Brokers, kafkaConfig.Topic, filters).Handle(handleEvent)
}

func handleEvent(c eventconsume.ConsumeContext) error {
	var event models.SaleRecordEvent
	if err := c.Bind(&event); err != nil {
		return err
	}

	ctx := c.Context()

	if err := (models.CustomerEventHandler{}).Handle(ctx, event); err != nil {
		return err
	}

	return nil
}
