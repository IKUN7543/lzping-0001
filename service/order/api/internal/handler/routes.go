package handler

import (
	"go-zero-ecommerce/service/order/api/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

func Routes(serverCtx *svc.ServiceContext) []rest.Route {
	return []rest.Route{
		{Method: "POST", Path: "/order/create", Handler: CreateOrderHandler(serverCtx)},
		{Method: "GET", Path: "/order/detail", Handler: GetOrderHandler(serverCtx)},
		{Method: "GET", Path: "/order/list", Handler: ListOrderHandler(serverCtx)},
		{Method: "POST", Path: "/order/cancel", Handler: CancelOrderHandler(serverCtx)},
		{Method: "POST", Path: "/order/pay", Handler: PayOrderHandler(serverCtx)},
	}
}
