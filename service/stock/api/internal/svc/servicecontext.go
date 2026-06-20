package svc

import (
	"go-zero-ecommerce/service/stock/api/internal/config"
	"go-zero-ecommerce/service/stock/rpc/stock"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	StockRpc stock.StockClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		StockRpc: stock.NewStockClient(zrpc.MustNewClient(c.StockRpc).Conn()),
	}
}
