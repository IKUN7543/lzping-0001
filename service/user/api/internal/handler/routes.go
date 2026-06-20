package handler

import (
	"go-zero-ecommerce/service/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func Routes(serverCtx *svc.ServiceContext) []rest.Route {
	return []rest.Route{
		{
			Method:  "POST",
			Path:    "/user/register",
			Handler: RegisterHandler(serverCtx),
		},
		{
			Method:  "POST",
			Path:    "/user/login",
			Handler: LoginHandler(serverCtx),
		},
		{
			Method:  "GET",
			Path:    "/user/v1/userinfo",
			Handler: UserInfoHandler(serverCtx),
		},
		{
			Method:  "PUT",
			Path:    "/user/v1/update",
			Handler: UpdateUserHandler(serverCtx),
		},
	}
}
