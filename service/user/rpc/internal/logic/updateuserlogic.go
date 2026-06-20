package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"go-zero-ecommerce/common/errx"
	"go-zero-ecommerce/service/user/rpc/internal/svc"
	"go-zero-ecommerce/service/user/rpc/model"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrUserNotFound
		}
		return nil, errx.ErrInternalServer
	}

	if in.Nickname != "" {
		u.Nickname = in.Nickname
	}
	if in.Avatar != "" {
		u.Avatar = in.Avatar
	}
	if in.Gender != 0 {
		u.Gender = in.Gender
	}

	err = l.svcCtx.UserModel.Update(l.ctx, u)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	return &user.UpdateUserResponse{
		User: &user.UserInfo{
			Id:       u.Id,
			Username: u.Username,
			Nickname: u.Nickname,
			Mobile:   u.Mobile,
			Email:    u.Email,
			Gender:   u.Gender,
			Avatar:   u.Avatar,
			Status:   u.Status,
		},
	}, nil
}

func toModelUser(u *user.UserInfo) *model.User {
	return &model.User{
		Id:       u.Id,
		Username: u.Username,
		Nickname: u.Nickname,
		Mobile:   u.Mobile,
		Email:    u.Email,
		Gender:   u.Gender,
		Avatar:   u.Avatar,
		Status:   u.Status,
	}
}
