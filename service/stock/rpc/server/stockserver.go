package server

import (
	"context"
	"go-zero-ecommerce/service/stock/rpc/internal/logic"
	"go-zero-ecommerce/service/stock/rpc/internal/svc"
	"go-zero-ecommerce/service/stock/rpc/stock"
)

type StockServer struct {
	svcCtx *svc.ServiceContext
	stock.UnimplementedStockServer
}

func NewStockServer(svcCtx *svc.ServiceContext) *StockServer {
	return &StockServer{svcCtx: svcCtx}
}

func (s *StockServer) CreateStock(ctx context.Context, in *stock.CreateStockRequest) (*stock.CreateStockResponse, error) {
	l := logic.NewCreateStockLogic(ctx, s.svcCtx)
	return l.CreateStock(in)
}

func (s *StockServer) GetStock(ctx context.Context, in *stock.GetStockRequest) (*stock.GetStockResponse, error) {
	l := logic.NewGetStockLogic(ctx, s.svcCtx)
	return l.GetStock(in)
}

func (s *StockServer) DeductStock(ctx context.Context, in *stock.DeductStockRequest) (*stock.DeductStockResponse, error) {
	l := logic.NewDeductStockLogic(ctx, s.svcCtx)
	return l.DeductStock(in)
}

func (s *StockServer) ReturnStock(ctx context.Context, in *stock.ReturnStockRequest) (*stock.ReturnStockResponse, error) {
	l := logic.NewReturnStockLogic(ctx, s.svcCtx)
	return l.ReturnStock(in)
}

func (s *StockServer) ConfirmStock(ctx context.Context, in *stock.ConfirmStockRequest) (*stock.ConfirmStockResponse, error) {
	l := logic.NewConfirmStockLogic(ctx, s.svcCtx)
	return l.ConfirmStock(in)
}
