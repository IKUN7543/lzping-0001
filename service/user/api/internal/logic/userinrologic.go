package logic

import (
	"context"
	"go-zero-ecommerce/service/user/api/internal/svc"
	"go-zero-ecommerce/service/user/api/internal/types"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	var userId int64
	if req.Id > 0 {
		userId = req.Id
	} else {
		uid, ok := l.ctx.Value("userId").(string)
		if !ok {
			userId = 0
		} else {
			var uidNum int64
			for _, c := range uid {
				uidNum = uidNum*10 + int64(c-'0')
			}
			userId = uidNum
		}
	}

	res, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserRequest{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResp{
		Id:        res.User.Id,
		Username:  res.User.Username,
		Nickname:  res.User.Nickname,
		Mobile:    res.User.Mobile,
		Email:     res.User.Email,
		Gender:    res.User.Gender,
		Avatar:    res.User.Avatar,
		Status:    res.User.Status,
		CreatedAt: res.User.CreatedAt,
		UpdatedAt: res.User.UpdatedAt,
	}, nil
}
