package logic

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"go-zero-ecommerce/common/errx"
	"go-zero-ecommerce/service/user/rpc/internal/svc"
	"go-zero-ecommerce/service/user/rpc/model"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	_, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err == nil {
		return nil, errx.ErrUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errx.ErrInternalServer
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	u := &model.User{
		Username: in.Username,
		Password: string(hashedPassword),
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Email:    in.Email,
		Status:   1,
	}
	if u.Nickname == "" {
		u.Nickname = u.Username
	}

	err = l.svcCtx.UserModel.Insert(l.ctx, u)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	return &user.RegisterResponse{
		User: &user.UserInfo{
			Id:       u.Id,
			Username: u.Username,
			Nickname: u.Nickname,
			Mobile:   u.Mobile,
			Email:    u.Email,
			Status:   u.Status,
		},
	}, nil
}
