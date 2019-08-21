package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nhub/sale-record-postprocess-api/adapters"
	"nhub/sale-record-postprocess-api/config"
	"nhub/sale-record-postprocess-api/factory"
	"nhub/sale-record-postprocess-api/models"
	"nomni/utils/auth"
	"nomni/utils/eventconsume"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

func main() {
	config := config.Init(os.Getenv("APP_ENV"))
	fmt.Println(config)
	saleRecordDB := initDB(config.Database.SaleRecord.Driver, config.Database.SaleRecord.Connection)
	orderDB := initDB(config.Database.Order.Driver, config.Database.Order.Connection)
	if err := models.InitSaleRecordDb(saleRecordDB); err != nil {
		panic(err)
	}
	if err := models.InitOrderDb(orderDB); err != nil {
		panic(err)
	}

	if config.SaleRecordEvent.Kafka.Brokers != nil {
		if err := adapters.NewSaleRecordEventConsumer(config.ServiceName, config.SaleRecordEvent.Kafka,
			eventconsume.Recover(),
			eventconsume.BehaviorLogger(config.ServiceName, config.BehaviorLog.Kafka),
			eventconsume.ContextDBWithName(config.ServiceName, factory.SaleRecordDBContextName, saleRecordDB, config.Database.Logger.Kafka),
			eventconsume.ContextDBWithName(config.ServiceName, factory.OrderDBContextName, orderDB, config.Database.Logger.Kafka),
		); err != nil {
			log.Fatal(err)
		}
	}

	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(echomiddleware.BehaviorLogger(config.ServiceName, config.BehaviorLog.Kafka))
	e.Use(echomiddleware.ContextDBWithName(config.ServiceName, echomiddleware.ContextDBType(factory.SaleRecordDBContextName), saleRecordDB, config.Database.Logger.Kafka))
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
