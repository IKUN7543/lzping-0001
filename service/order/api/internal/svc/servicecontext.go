package svc

import (
	"go-zero-ecommerce/service/order/api/internal/config"
	"go-zero-ecommerce/service/order/rpc/order"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	OrderRpc order.OrderClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		OrderRpc: order.NewOrderClient(zrpc.MustNewClient(c.OrderRpc).Conn()),
	}
}
