package svc

import (
	"go-zero-ecommerce/service/user/api/internal/config"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
	}
}
