package logic

import (
	"context"
	"go-zero-ecommerce/service/user/api/internal/svc"
	"go-zero-ecommerce/service/user/api/internal/types"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) (resp *types.UpdateUserResp, err error) {
	uid, ok := l.ctx.Value("userId").(string)
	var userId int64
	if ok {
		for _, c := range uid {
			userId = userId*10 + int64(c-'0')
		}
	}

	res, err := l.svcCtx.UserRpc.UpdateUser(l.ctx, &user.UpdateUserRequest{
		Id:       userId,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateUserResp{
		Id:       res.User.Id,
		Nickname: res.User.Nickname,
		Avatar:   res.User.Avatar,
		Gender:   res.User.Gender,
	}, nil
}
