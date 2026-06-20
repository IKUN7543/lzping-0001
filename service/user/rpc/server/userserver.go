package server

import (
	"context"
	"go-zero-ecommerce/service/user/rpc/internal/logic"
	"go-zero-ecommerce/service/user/rpc/internal/svc"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Register(ctx context.Context, in *user.RegisterRequest) (*user.RegisterResponse, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	logx.Infof("Register request: username=%s", in.Username)
	return l.Register(in)
}

func (s *UserServer) Login(ctx context.Context, in *user.LoginRequest) (*user.LoginResponse, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	logx.Infof("Login request: username=%s", in.Username)
	return l.Login(in)
}

func (s *UserServer) GetUser(ctx context.Context, in *user.GetUserRequest) (*user.GetUserResponse, error) {
	l := logic.NewGetUserLogic(ctx, s.svcCtx)
	return l.GetUser(in)
}

func (s *UserServer) UpdateUser(ctx context.Context, in *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	l := logic.NewUpdateUserLogic(ctx, s.svcCtx)
	return l.UpdateUser(in)
}
