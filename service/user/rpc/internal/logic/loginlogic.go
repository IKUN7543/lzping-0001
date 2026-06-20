package logic

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"go-zero-ecommerce/common/errx"
	"go-zero-ecommerce/service/user/rpc/internal/svc"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrUserNotFound
		}
		return nil, errx.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password))
	if err != nil {
		return nil, errx.ErrPasswordWrong
	}

	if u.Status != 1 {
		return nil, errx.NewError(20004, "user account is disabled")
	}

	return &user.LoginResponse{
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
