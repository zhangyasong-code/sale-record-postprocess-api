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
	"nhub/sale-record-postprocess-api/promotion"

	"nhub/sale-record-postprocess-api/factory"
	"nhub/sale-record-postprocess-api/models"
	"nomni/utils/auth"
	"nomni/utils/eventconsume"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
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

	if err := models.InitSaleRecordDb(saleRecordDB); err != nil {
		log.Fatal(err)
	}
	if err := models.InitOrderDb(orderDB); err != nil {
		log.Fatal(err)
	}
	if err := promotion.InitDB(db); err != nil {
		log.Fatal(err)
	}

	if err := adapters.NewConsumers(config.ServiceName, config.Kafka,
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

	controllers.PromotionEventController{}.Init(r.Group("PromotionEvent", "/v1/promotion-event"))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(echomiddleware.ContextBase())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDBWithName(config.ServiceName, factory.SaleRecordDBContextName, saleRecordDB, echomiddleware.KafkaConfig(config.Database.Logger.Kafka)))
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
