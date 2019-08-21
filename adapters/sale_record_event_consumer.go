package adapters

import (
	"nomni/utils/eventconsume"

	"github.com/pangpanglabs/goutils/kafka"
)

func NewSaleRecordEventConsumer(serviceName string, kafkaConfig kafka.Config, filters ...eventconsume.Filter) error {
	return eventconsume.NewEventConsumer(serviceName, kafkaConfig.Brokers, kafkaConfig.Topic, filters).Handle(handleEvent)
}

func handleEvent(c eventconsume.ConsumeContext) error {
	return nil
}
