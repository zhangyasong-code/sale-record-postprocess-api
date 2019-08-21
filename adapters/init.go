package adapters

import (
	"errors"
	"nhub/sale-record-postprocess-api/config"
	"nomni/utils/eventconsume"

	"github.com/sirupsen/logrus"
)

func NewConsumers(serviceName string, kafkaConfig config.EventKafka, filters ...eventconsume.Filter) error {
	saleErr := NewSaleRecordEventConsumer(serviceName, kafkaConfig.SaleRecordEvent, filters...)
	if saleErr != nil {
		logrus.WithFields(logrus.Fields{
			"err": saleErr,
		}).Error("NewSaleRecordEventConsumer Error")
	}
	promotionErr := NewPromotionEventConsumer(serviceName, kafkaConfig.PromotionEvent, filters...)
	if promotionErr != nil {
		logrus.WithFields(logrus.Fields{
			"err": promotionErr,
		}).Error("NewPromotionEventConsumer Error")
	}

	if promotionErr != nil && saleErr != nil {
		return errors.New("NewConsumer Error")
	}
	return nil
}
