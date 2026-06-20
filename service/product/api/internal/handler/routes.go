package handler

import (
	"go-zero-ecommerce/service/product/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func Routes(serverCtx *svc.ServiceContext) []rest.Route {
	return []rest.Route{
		{
			Method:  "GET",
			Path:    "/product/detail",
			Handler: ProductInfoHandler(serverCtx),
		},
		{
			Method:  "GET",
			Path:    "/product/list",
			Handler: ListProductHandler(serverCtx),
		},
		{
			Method:  "GET",
			Path:    "/product/search",
			Handler: SearchProductHandler(serverCtx),
		},
		{
			Method:  "GET",
			Path:    "/product/category",
			Handler: ListCategoryHandler(serverCtx),
		},
		{
			Method:  "POST",
			Path:    "/product/admin/create",
			Handler: CreateProductHandler(serverCtx),
		},
		{
			Method:  "POST",
			Path:    "/product/admin/update",
			Handler: UpdateProductHandler(serverCtx),
		},
	}
}
