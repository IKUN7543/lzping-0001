package svc

import (
	"go-zero-ecommerce/service/user/rpc/config"
	"go-zero-ecommerce/service/user/rpc/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(db),
	}
}
