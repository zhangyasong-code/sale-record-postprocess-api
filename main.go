package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nhub/sale-record-postprocess-api/adapters"
	"nhub/sale-record-postprocess-api/config"
	"nhub/sale-record-postprocess-api/controllers"
	"nhub/sale-record-postprocess-api/customer"
	"nhub/sale-record-postprocess-api/promotion"
	"nhub/sale-record-postprocess-api/salePerson"
	"nhub/sale-record-postprocess-api/saleRecordFee"

	"nhub/sale-record-postprocess-api/factory"
	"nomni/utils/auth"
	"nomni/utils/eventconsume"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pangpanglabs/echoswagger"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

func main() {
	config := config.Init(os.Getenv("APP_ENV"))
	fmt.Println(config)
	saleRecordDB := initDB(config.Database.SaleRecord.Driver, config.Database.SaleRecord.Connection)
	orderDB := initDB(config.Database.Order.Driver, config.Database.Order.Connection)
	defer saleRecordDB.Close()
	defer orderDB.Close()

	if err := customer.InitDB(saleRecordDB); err != nil {
		log.Fatal(err)
	}
	if err := payamt.InitDB(saleRecordDB); err != nil {
		log.Fatal(err)
	}
	if err := promotion.InitDB(saleRecordDB); err != nil {
		log.Fatal(err)
	}
	if err := salePerson.InitDB(saleRecordDB); err != nil {
		log.Fatal(err)
	}

	if err := saleRecordFee.InitDB(saleRecordDB); err != nil {
		log.Fatal(err)
	}

	if err := adapters.NewConsumers(config.ServiceName, config.EventKafka,
		eventconsume.Recover(),
		eventconsume.BehaviorLogger(config.ServiceName, config.BehaviorLog.Kafka),
		eventconsume.ContextDBWithName(config.ServiceName, factory.SaleRecordDBContextName, saleRecordDB, config.Database.Logger.Kafka),
		eventconsume.ContextDBWithName(config.ServiceName, factory.OrderDBContextName, orderDB, config.Database.Logger.Kafka),
	); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	r := echoswagger.New(e, "docs", &echoswagger.Info{
		Title:       "Sale Record Postprocess API",
		Description: "This is docs for sale-record-postprocess-api service",
		Version:     "1.0.0",
	})

	r.AddSecurityAPIKey("Authorization", "JWT token", echoswagger.SecurityInHeader)
	r.SetUI(echoswagger.UISetting{
		HideTop: true,
	})

	controllers.SaleRecordEventController{}.Init(r.Group("SaleRecordEvent", "/v1/saleRecord-events"))
	controllers.PromotionEventController{}.Init(r.Group("PromotionEvent", "/v1/promotion-event"))
	controllers.SaleRecordInfoController{}.Init(r.Group("SaleRecordInfo", "/v1/sale-record-info"))
	salePerson.SalesPersonEventHandler{}.Init(r.Group("SalesPerson", "/v1/sales-person"))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(echomiddleware.ContextBase())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDBWithName(config.ServiceName, echomiddleware.ContextDBType(factory.SaleRecordDBContextName), saleRecordDB, echomiddleware.KafkaConfig(config.Database.Logger.Kafka)))
	e.Use(echomiddleware.BehaviorLogger(config.ServiceName, config.BehaviorLog.Kafka))
	e.Use(auth.UserClaimMiddleware("/ping", "/docs"))

	if err := e.Start(":8000"); err != nil {
		log.Println(err)
	}
}

func initDB(driver, connection string) *xorm.Engine {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		panic(err)
	}
	if os.Getenv("APP_ENV") != "production" {
		db.ShowSQL(true)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Minute * 10)

	return db
}
