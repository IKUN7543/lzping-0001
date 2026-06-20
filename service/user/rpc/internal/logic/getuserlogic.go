package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"go-zero-ecommerce/common/errx"
	"go-zero-ecommerce/service/user/rpc/internal/svc"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserRequest) (*user.GetUserResponse, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrUserNotFound
		}
		return nil, errx.ErrInternalServer
	}

	return &user.GetUserResponse{
		User: &user.UserInfo{
			Id:        u.Id,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Mobile:    u.Mobile,
			Email:     u.Email,
			Gender:    u.Gender,
			Avatar:    u.Avatar,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
