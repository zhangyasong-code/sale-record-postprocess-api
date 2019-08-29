package adapters

import (
	"encoding/json"
	"nhub/sale-record-postprocess-api/models"
	"nomni/utils/eventconsume"

	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

func NewSaleRecordEventConsumer(serviceName string, kafkaConfig kafka.Config, filters ...eventconsume.Filter) error {
	return eventconsume.NewEventConsumer(serviceName, kafkaConfig.Brokers, kafkaConfig.Topic, filters).Handle(handleEvent)
}

func handleEvent(c eventconsume.ConsumeContext) error {
	var event models.SaleRecordEvent
	if err := c.Bind(&event); err != nil {
		logrus.WithField("Error", err).Info("Event bind error!")
		return err
	}

	ctx := c.Context()
	str, _ := json.Marshal(event)
	logrus.WithField("Body", string(str)).Info("Event Body>>>>>>")
	if err := (models.CustomerEventHandler{}).Handle(ctx, event); err != nil {
		return err
	}

	return nil
}
