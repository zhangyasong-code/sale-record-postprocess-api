package adapters

import (
	"encoding/json"
	"nhub/sale-record-postprocess-api/customer"
	"nhub/sale-record-postprocess-api/models"
	"nhub/sale-record-postprocess-api/payamt"
	"nhub/sale-record-postprocess-api/salePerson"
	"nhub/sale-record-postprocess-api/saleRecordFee"
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

	if err := (payamt.PayAmtEventHandler{}).Handle(ctx, event); err != nil {
		return err
	}
	if err := (customer.CustomerEventHandler{}).Handle(ctx, event); err != nil {
		return err
	}

	if err := (salePerson.SalesPersonEventHandler{}).Handle(ctx, event); err != nil {
		return err
	}

	if err := (saleRecordFee.SaleRecordFeeEventHandler{}).Handle(ctx, event); err != nil {
		return err
	}

	return nil
}
