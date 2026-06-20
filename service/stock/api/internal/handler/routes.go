package handler

import (
	"go-zero-ecommerce/service/stock/api/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

func Routes(serverCtx *svc.ServiceContext) []rest.Route {
	return []rest.Route{
		{Method: "GET", Path: "/stock/detail", Handler: GetStockHandler(serverCtx)},
		{Method: "POST", Path: "/stock/admin/create", Handler: CreateStockHandler(serverCtx)},
		{Method: "POST", Path: "/stock/deduct", Handler: DeductStockHandler(serverCtx)},
	}
}
