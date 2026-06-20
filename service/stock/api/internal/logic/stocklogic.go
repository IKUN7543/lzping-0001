package logic

import (
	"context"
	"go-zero-ecommerce/service/stock/api/internal/svc"
	"go-zero-ecommerce/service/stock/api/internal/types"
	stock2 "go-zero-ecommerce/service/stock/rpc/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetStockLogic) GetStock(req *types.GetStockReq) (resp *types.GetStockResp, err error) {
	res, err := l.svcCtx.StockRpc.GetStock(l.ctx, &stock2.GetStockRequest{ProductId: req.ProductId})
	if err != nil {
		return nil, err
	}
	return &types.GetStockResp{
		Id: res.Stock.Id, ProductId: res.Stock.ProductId, Total: res.Stock.Total,
		Available: res.Stock.Available, LockStock: res.Stock.LockStock, Sales: res.Stock.Sales,
	}, nil
}

type CreateStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStockLogic {
	return &CreateStockLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateStockLogic) CreateStock(req *types.CreateStockReq) (resp *types.CreateStockResp, err error) {
	res, err := l.svcCtx.StockRpc.CreateStock(l.ctx, &stock2.CreateStockRequest{ProductId: req.ProductId, Num: req.Num})
	if err != nil {
		return nil, err
	}
	return &types.CreateStockResp{
		ProductId: res.Stock.ProductId, Total: res.Stock.Total, Available: res.Stock.Available,
	}, nil
}

type DeductStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockLogic {
	return &DeductStockLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeductStockLogic) DeductStock(req *types.DeductStockReq) (resp *types.DeductStockResp, err error) {
	res, err := l.svcCtx.StockRpc.DeductStock(l.ctx, &stock2.DeductStockRequest{
		ProductId: req.ProductId, Num: req.Num, OrderNo: req.OrderNo,
	})
	if err != nil {
		return &types.DeductStockResp{Success: false}, nil
	}
	return &types.DeductStockResp{Success: res.Success}, nil
}
