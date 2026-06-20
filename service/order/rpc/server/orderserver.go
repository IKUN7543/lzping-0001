package server

import (
	"context"
	"go-zero-ecommerce/service/order/rpc/internal/logic"
	"go-zero-ecommerce/service/order/rpc/internal/svc"
	"go-zero-ecommerce/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderServer struct {
	svcCtx *svc.ServiceContext
	order.UnimplementedOrderServer
}

func NewOrderServer(svcCtx *svc.ServiceContext) *OrderServer {
	return &OrderServer{svcCtx: svcCtx}
}

func (s *OrderServer) CreateOrder(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	l := logic.NewCreateOrderLogic(ctx, s.svcCtx)
	logx.Infof("CreateOrder request: userId=%d, items=%d", in.UserId, len(in.Items))
	return l.CreateOrder(in)
}

func (s *OrderServer) GetOrder(ctx context.Context, in *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	l := logic.NewGetOrderLogic(ctx, s.svcCtx)
	return l.GetOrder(in)
}

func (s *OrderServer) ListOrder(ctx context.Context, in *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	l := logic.NewListOrderLogic(ctx, s.svcCtx)
	return l.ListOrder(in)
}

func (s *OrderServer) CancelOrder(ctx context.Context, in *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	l := logic.NewCancelOrderLogic(ctx, s.svcCtx)
	logx.Infof("CancelOrder request: orderNo=%s", in.OrderNo)
	return l.CancelOrder(in)
}

func (s *OrderServer) PayOrder(ctx context.Context, in *order.PayOrderRequest) (*order.PayOrderResponse, error) {
	l := logic.NewPayOrderLogic(ctx, s.svcCtx)
	logx.Infof("PayOrder request: orderNo=%s", in.OrderNo)
	return l.PayOrder(in)
}
