package svc

import (
	"go-zero-ecommerce/service/product/api/internal/config"
	"go-zero-ecommerce/service/product/rpc/product"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	ProductRpc product.ProductClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ProductRpc: product.NewProductClient(zrpc.MustNewClient(c.ProductRpc).Conn()),
	}
}
