package adapters

import (
	"encoding/json"
	"fmt"
	"nhub/sale-record-postprocess-api/customer"
	"nhub/sale-record-postprocess-api/models"
	"nhub/sale-record-postprocess-api/payamt"
	"nhub/sale-record-postprocess-api/postprocess"
	"nhub/sale-record-postprocess-api/refundApproval"
	"nhub/sale-record-postprocess-api/salePerson"
	"nhub/sale-record-postprocess-api/saleRecordFee"
	"nhub/sale-record-postprocess-api/sendCsl"
	"nomni/utils/eventconsume"
	"strings"

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

	// TransactionChannelType > EMALL AND  TransactionType > EXCHANGE not to refundApproval
	if event.RefundId > 0 && strings.ToUpper(event.TransactionChannelType) != "EMALL" && strings.ToUpper(event.TransactionType) != "EXCHANGE" {
		isAllowTransCSL, err := refundApproval.Check(ctx, event.TenantCode, event.StoreId, event.OrderId, event.RefundId, event.Committed.Created)
		if err != nil {
			postProcessSuccess := &postprocess.PostProcessSuccess{
				TransactionId: event.TransactionId,
				OrderId:       event.OrderId,
				RefundId:      event.RefundId,
				StoreId:       event.StoreId,
				ModuleType:    string(postprocess.ModuleRefundApproval),
				IsSuccess:     false,
				Error:         err.Error(),
				Details:       postprocess.GetErrorDetails(string(postprocess.ModuleRefundApproval)),
				ModuleEntity:  string(str),
			}
			if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
				return saveErr
			}
			return fmt.Errorf("Refund Approval Error")
		} else if !isAllowTransCSL {
			postProcessSuccess := &postprocess.PostProcessSuccess{
				TransactionId: event.TransactionId,
				OrderId:       event.OrderId,
				RefundId:      event.RefundId,
				StoreId:       event.StoreId,
				ModuleType:    string(postprocess.ModuleRefundApproval),
				IsSuccess:     false,
				Error:         "Refund Approval Error",
				Details:       "退货审批中！",
				ModuleEntity:  string(str),
			}
			if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
				return saveErr
			}
			return fmt.Errorf("Refund Approval Error")
		} else {
			postProcessSuccess := &postprocess.PostProcessSuccess{
				TransactionId: event.TransactionId,
				OrderId:       event.OrderId,
				RefundId:      event.RefundId,
				StoreId:       event.StoreId,
				ModuleType:    string(postprocess.ModuleRefundApproval),
				IsSuccess:     true,
				Error:         "",
				ModuleEntity:  string(str),
			}
			if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
				return saveErr
			}
		}
	}

	if err := (payamt.PayAmtEventHandler{}).Handle(ctx, event); err != nil {
		logrus.WithFields(logrus.Fields{
			"TransactionId": event.TransactionId,
			"OrderId":       event.OrderId,
		}).WithError(err).Error("Fail to handle PayAmtEventHandler")

		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModulePay),
			IsSuccess:     false,
			Error:         err.Error(),
			Details:       postprocess.GetErrorDetails(string(postprocess.ModulePay)),
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}

		return err
	} else {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModulePay),
			IsSuccess:     true,
			Error:         "",
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}
	}

	if err := (customer.CustomerEventHandler{}).Handle(ctx, event); err != nil {
		logrus.WithFields(logrus.Fields{
			"TransactionId": event.TransactionId,
			"OrderId":       event.OrderId,
			"RefundId":      event.RefundId,
		}).WithError(err).Error("Fail to handle CustomerEventHandler")

		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModuleMileage),
			IsSuccess:     false,
			Error:         err.Error(),
			Details:       postprocess.GetErrorDetails(string(postprocess.ModuleMileage)),
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}

		return err
	} else {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModuleMileage),
			IsSuccess:     true,
			Error:         "",
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}
	}

	if err := (salePerson.SalesPersonEventHandler{}).Handle(ctx, event); err != nil {
		logrus.WithFields(logrus.Fields{
			"TransactionId": event.TransactionId,
			"OrderId":       event.OrderId,
			"RefundId":      event.RefundId,
		}).WithError(err).Error("Fail to handle SalesPersonEventHandler")

		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModuleSalePerson),
			IsSuccess:     false,
			Error:         err.Error(),
			Details:       postprocess.GetErrorDetails(string(postprocess.ModuleSalePerson)),
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}

		return err
	} else {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModuleSalePerson),
			IsSuccess:     true,
			Error:         "",
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}
	}

	if err := (saleRecordFee.SaleRecordFeeEventHandler{}).Handle(ctx, event); err != nil {
		logrus.WithFields(logrus.Fields{
			"TransactionId": event.TransactionId,
			"OrderId":       event.OrderId,
			"RefundId":      event.RefundId,
		}).WithError(err).Error("Fail to handle SaleRecordFeeEventHandler")

		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModuleSaleFee),
			IsSuccess:     false,
			Error:         err.Error(),
			Details:       postprocess.GetErrorDetails(string(postprocess.ModuleSaleFee)),
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}

		return err
	} else {
		postProcessSuccess := &postprocess.PostProcessSuccess{
			TransactionId: event.TransactionId,
			OrderId:       event.OrderId,
			RefundId:      event.RefundId,
			StoreId:       event.StoreId,
			ModuleType:    string(postprocess.ModuleSaleFee),
			IsSuccess:     true,
			Error:         "",
			ModuleEntity:  string(str),
		}
		if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
			return saveErr
		}
	}
	logrus.WithFields(logrus.Fields{
		"TransactionId": event.TransactionId,
		"OrderId":       event.OrderId,
		"RefundId":      event.RefundId,
	}).Info("Success to handle event")

	// Send to Csl
	if event.TransactionStatus == "BuyerReceivedConfirmed" || event.TransactionStatus == "RefundOrderSuccess" {
		if err := (sendCsl.Send{}).SendToCsl(ctx, event); err != nil {
			logrus.WithFields(logrus.Fields{
				"TransactionId": event.TransactionId,
				"OrderId":       event.OrderId,
				"RefundId":      event.RefundId,
			}).WithError(err).Error("Fail to SendToCsl")

			postProcessSuccess := &postprocess.PostProcessSuccess{
				TransactionId: event.TransactionId,
				OrderId:       event.OrderId,
				RefundId:      event.RefundId,
				StoreId:       event.StoreId,
				ModuleType:    string(postprocess.SendToClearance),
				IsSuccess:     false,
				Error:         err.Error(),
				Details:       postprocess.GetErrorDetails(string(postprocess.SendToClearance)),
				ModuleEntity:  string(str),
			}
			if saveErr := postProcessSuccess.Save(ctx); saveErr != nil {
				return saveErr
			}
			return err
		}
		postProcessSuccess, err := postprocess.PostProcessSuccess{}.Get(ctx, false, event.TransactionId, string(postprocess.SendToClearance))
		if err != nil {
			return err
		}
		if postProcessSuccess.Id != 0 {
			postProcessSuccess.IsSuccess = true
			if err := postProcessSuccess.Update(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}
