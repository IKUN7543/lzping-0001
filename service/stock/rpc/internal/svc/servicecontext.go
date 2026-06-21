package svc

import (
	"go-zero-ecommerce/service/stock/rpc/config"
	"go-zero-ecommerce/service/stock/rpc/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	StockModel model.StockModel
	Rdb        *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	poolSize := c.Redis.PoolSize
	if poolSize <= 0 {
		poolSize = 100
	}
	minIdle := c.Redis.MinIdleConns
	if minIdle <= 0 {
		minIdle = 10
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Host,
		Password:     c.Redis.Pass,
		DB:           0,
		PoolSize:     poolSize,
		MinIdleConns: minIdle,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
		MaxRetries:   3,
	})

	return &ServiceContext{
		Config:     c,
		StockModel: model.NewStockModel(db),
		Rdb:        rdb,
	}
}
