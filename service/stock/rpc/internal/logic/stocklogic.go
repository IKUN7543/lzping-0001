package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-zero-ecommerce/common/errx"
	redislock2 "go-zero-ecommerce/common/redislock"
	"go-zero-ecommerce/service/stock/rpc/internal/svc"
	"go-zero-ecommerce/service/stock/rpc/model"
	"go-zero-ecommerce/service/stock/rpc/stock"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStockLogic {
	return &CreateStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateStockLogic) CreateStock(in *stock.CreateStockRequest) (*stock.CreateStockResponse, error) {
	s := &model.Stock{
		ProductId: in.ProductId,
		Total:   in.Num,
		Available: in.Num,
	}
	err := l.svcCtx.StockModel.Insert(l.ctx, s)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	stockKey := fmt.Sprintf("stock:%d", in.ProductId)
	l.svcCtx.Rdb.Set(l.ctx, stockKey, in.Num, 0)

	return &stock.CreateStockResponse{Stock: &stock.StockInfo{
		Id: s.Id, ProductId: s.ProductId, Total: s.Total, Available: s.Available,
	}}, nil
}

type GetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetStockLogic) GetStock(in *stock.GetStockRequest) (*stock.GetStockResponse, error) {
	s, err := l.svcCtx.StockModel.FindByProductId(l.ctx, in.ProductId)
	if err != nil {
		return nil, errx.ErrStockNotFound
	}
	return &stock.GetStockResponse{Stock: &stock.StockInfo{
		Id: s.Id, ProductId: s.ProductId, Total: s.Total, Available: s.Available, LockStock: s.LockStock, Sales: s.Sales,
	}}, nil
}

type DeductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockLogic {
	return &DeductStockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeductStockLogic) DeductStock(in *stock.DeductStockRequest) (*stock.DeductStockResponse, error) {
	if in.Num <= 0 {
		return nil, errx.ErrInvalidParam
	}

	deductKey := fmt.Sprintf("stock_deduct:%s", in.OrderNo)
	exists, _ := l.svcCtx.Rdb.Exists(l.ctx, deductKey).Result()
	if exists > 0 {
		return &stock.DeductStockResponse{Success: true}, nil
	}

	lockKey := fmt.Sprintf("stock_lock:%d", in.ProductId)
	rlock := redislock2.NewRedisLock(l.svcCtx.Rdb, lockKey, 10*time.Second)
	err := rlock.TryLock(l.ctx, 5, 100*time.Millisecond)
	if err != nil {
		return &stock.DeductStockResponse{Success: false}, errx.ErrStockNotEnough
	}
	defer rlock.Unlock(l.ctx)

	redisStockKey := fmt.Sprintf("stock:%d", in.ProductId)
	available, err := l.svcCtx.Rdb.DecrBy(l.ctx, redisStockKey, int64(in.Num)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return &stock.DeductStockResponse{Success: false}, errx.ErrInternalServer
	}
	if available < 0 {
		l.svcCtx.Rdb.IncrBy(l.ctx, redisStockKey, int64(in.Num))
		return &stock.DeductStockResponse{Success: false}, errx.ErrStockNotEnough
	}

	s, err := l.svcCtx.StockModel.FindByProductId(l.ctx, in.ProductId)
	if err != nil {
		l.svcCtx.Rdb.IncrBy(l.ctx, redisStockKey, int64(in.Num))
		return &stock.DeductStockResponse{Success: false}, errx.ErrStockNotFound
	}

	affected, err := l.svcCtx.StockModel.DecrAvailable(l.ctx, in.ProductId, in.Num, s.Version)
	if err != nil || affected == 0 {
		l.svcCtx.Rdb.IncrBy(l.ctx, redisStockKey, int64(in.Num))
		return &stock.DeductStockResponse{Success: false}, errx.ErrStockNotEnough
	}

	l.svcCtx.Rdb.SetNX(l.ctx, deductKey, "1", 24*time.Hour)
	logx.Infof("Deduct stock success: productId=%d, num=%d, orderNo=%s", in.ProductId, in.Num, in.OrderNo)
	return &stock.DeductStockResponse{Success: true}, nil
}

type ReturnStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReturnStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReturnStockLogic {
	return &ReturnStockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReturnStockLogic) ReturnStock(in *stock.ReturnStockRequest) (*stock.ReturnStockResponse, error) {
	affected, err := l.svcCtx.StockModel.IncrAvailable(l.ctx, in.ProductId, in.Num)
	if err != nil || affected == 0 {
		return &stock.ReturnStockResponse{Success: false}, errx.ErrInternalServer
	}

	redisStockKey := fmt.Sprintf("stock:%d", in.ProductId)
	l.svcCtx.Rdb.IncrBy(l.ctx, redisStockKey, int64(in.Num))

	return &stock.ReturnStockResponse{Success: true}, nil
}

type ConfirmStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmStockLogic {
	return &ConfirmStockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConfirmStockLogic) ConfirmStock(in *stock.ConfirmStockRequest) (*stock.ConfirmStockResponse, error) {
	s, err := l.svcCtx.StockModel.FindByProductId(l.ctx, in.ProductId)
	if err != nil {
		return &stock.ConfirmStockResponse{Success: false}, errx.ErrStockNotFound
	}
	s.LockStock -= in.Num
	s.Sales += in.Num
	err = l.svcCtx.StockModel.Update(l.ctx, s)
	if err != nil {
		return &stock.ConfirmStockResponse{Success: false}, errx.ErrInternalServer
	}

	confirmKey := "stock_confirm:" + in.OrderNo
	l.svcCtx.Rdb.Set(l.ctx, confirmKey, strconv.Itoa(int(in.Num)), 24*time.Hour)
	return &stock.ConfirmStockResponse{Success: true}, nil
}
