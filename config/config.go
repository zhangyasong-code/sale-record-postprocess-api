package config

import (
	"log"
	"os"

	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/jwtutil"
	"github.com/sirupsen/logrus"
)

var config C

func Init(appEnv string, options ...func(*C)) C {
	config.AppEnv = appEnv
	if err := configutil.Read(appEnv, &config); err != nil {
		logrus.WithError(err).Warn("Fail to load config file")
	}

	log.Println("APP_ENV:", appEnv)
	log.Printf("config: %+v\n", config)

	if s := os.Getenv("JWT_SECRET"); s != "" {
		config.JwtSecret = s
		jwtutil.SetJwtSecret(s)
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

func Config() C {
	return config
}

type EventKafka struct {
	SaleRecordEvent echomiddleware.KafkaConfig
	PromotionEvent  echomiddleware.KafkaConfig
}

type C struct {
	Database struct {
		SaleRecord struct {
			Driver     string
			Connection string
		}
		Order struct {
			Driver     string
			Connection string
		}
		Logger struct {
			Kafka echomiddleware.KafkaConfig
		}
	}
	BehaviorLog struct {
		Kafka echomiddleware.KafkaConfig
	}
	EventKafka EventKafka

	StockSourceTypes []string
	Services         struct {
		BenefitApi,
		ProductApi,
		OfferApi,
		PromotionApi,
		StoreApi,
		PlaceManagementApi string
	}
	AppEnv, JwtSecret string
	ServiceName       string
	Debug             bool
}
