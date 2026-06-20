package logic

import (
	"context"
	"strconv"
	"time"
	"go-zero-ecommerce/common/errx"
	"go-zero-ecommerce/service/user/api/internal/svc"
	"go-zero-ecommerce/service/user/api/internal/types"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, errx.ErrInvalidParam
	}

	res, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Mobile:   req.Mobile,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		Id:       res.User.Id,
		Username: res.User.Username,
		Nickname: res.User.Nickname,
	}, nil
}

func (l *RegisterLogic) genToken(userId int64) (accessToken string, refreshToken string, expiresAt int64, err error) {
	secret := l.svcCtx.Config.JwtAuth.AccessSecret
	expire := l.svcCtx.Config.JwtAuth.AccessExpire
	now := time.Now().Unix()

	accessClaims := make(jwt.MapClaims)
	accessClaims["exp"] = now + expire
	accessClaims["iat"] = now
	accessClaims["userId"] = strconv.FormatInt(userId, 10)
	accessTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenJwt.SignedString([]byte(secret))
	if err != nil {
		return
	}

	refreshExpire := expire * 7
	refreshClaims := make(jwt.MapClaims)
	refreshClaims["exp"] = now + refreshExpire
	refreshClaims["iat"] = now
	refreshClaims["userId"] = strconv.FormatInt(userId, 10)
	refreshTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenJwt.SignedString([]byte(secret))
	if err != nil {
		return
	}

	expiresAt = now + expire
	return
}
