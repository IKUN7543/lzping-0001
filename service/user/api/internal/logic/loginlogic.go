package logic

import (
	"context"
	"strconv"
	"time"
	"go-zero-ecommerce/service/user/api/internal/svc"
	"go-zero-ecommerce/service/user/api/internal/types"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, expiresAt, err := l.genToken(res.User.Id)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Id:           res.User.Id,
		Username:     res.User.Username,
		Nickname:     res.User.Nickname,
	}, nil
}

func (l *LoginLogic) genToken(userId int64) (accessToken string, refreshToken string, expiresAt int64, err error) {
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
