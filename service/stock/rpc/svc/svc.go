package svc

import (
	"github.com/redis/go-redis/v9"
	"go-zero-ecommerce/service/stock/rpc/config"
	"go-zero-ecommerce/service/stock/rpc/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	StockModel model.StockModel
	Rdb       *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       0,
	})

	return &ServiceContext{
		Config:    c,
		StockModel: model.NewStockModel(db),
		Rdb:       rdb,
	}
}
