package logic

import (
	"context"
	"go-zero-ecommerce/service/order/api/internal/svc"
	"go-zero-ecommerce/service/order/api/internal/types"
	order2 "go-zero-ecommerce/service/order/rpc/order"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

func getUserIdFromCtx(ctx context.Context) int64 {
	v := ctx.Value("userId")
	if v == nil {
		return 0
	}
	s, ok := v.(string)
	if !ok {
		return 0
	}
	id, _ := strconv.ParseInt(s, 10, 64)
	return id
}

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {
	userId := getUserIdFromCtx(l.ctx)
	items := make([]*order2.OrderItemReq, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, &order2.OrderItemReq{
			ProductId: it.ProductId, Num: it.Num, Price: it.Price,
		})
	}

	res, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order2.CreateOrderRequest{
		UserId:          userId,
		Items:           items,
		ReceiverName:    req.ReceiverName,
		ReceiverPhone:   req.ReceiverPhone,
		ReceiverAddress: req.ReceiverAddress,
		Remark:          req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateOrderResp{
		OrderNo: res.OrderNo, Id: res.Id, PayAmount: res.PayAmount,
	}, nil
}

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderReq) (resp *types.GetOrderResp, err error) {
	res, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order2.GetOrderRequest{OrderNo: req.OrderNo})
	if err != nil {
		return nil, err
	}

	items := make([]types.OrderItemResp, 0, len(res.Order.Items))
	for _, it := range res.Order.Items {
		items = append(items, types.OrderItemResp{
			Id: it.Id, ProductId: it.ProductId, ProductName: it.ProductName,
			ProductImage: it.ProductImage, Price: it.Price, Num: it.Num, TotalPrice: it.TotalPrice,
		})
	}

	return &types.GetOrderResp{Order: types.OrderResp{
		Id: res.Order.Id, OrderNo: res.Order.OrderNo, UserId: res.Order.UserId,
		TotalAmount: res.Order.TotalAmount, PayAmount: res.Order.PayAmount,
		FreightAmount: res.Order.FreightAmount, DiscountAmount: res.Order.DiscountAmount,
		Status: res.Order.Status, PayType: res.Order.PayType,
		ReceiverName: res.Order.ReceiverName, ReceiverPhone: res.Order.ReceiverPhone,
		ReceiverAddress: res.Order.ReceiverAddress, Remark: res.Order.Remark,
		Items: items, CreatedAt: res.Order.CreatedAt, UpdatedAt: res.Order.UpdatedAt,
	}}, nil
}

type ListOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderLogic {
	return &ListOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListOrderLogic) ListOrder(req *types.ListOrderReq) (resp *types.ListOrderResp, err error) {
	userId := getUserIdFromCtx(l.ctx)
	res, err := l.svcCtx.OrderRpc.ListOrder(l.ctx, &order2.ListOrderRequest{
		UserId: userId, Page: req.Page, PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]types.OrderResp, 0, len(res.Orders))
	for _, o := range res.Orders {
		items := make([]types.OrderItemResp, 0, len(o.Items))
		for _, it := range o.Items {
			items = append(items, types.OrderItemResp{
				Id: it.Id, ProductId: it.ProductId, ProductName: it.ProductName,
				ProductImage: it.ProductImage, Price: it.Price, Num: it.Num, TotalPrice: it.TotalPrice,
			})
		}
		orders = append(orders, types.OrderResp{
			Id: o.Id, OrderNo: o.OrderNo, UserId: o.UserId,
			TotalAmount: o.TotalAmount, PayAmount: o.PayAmount, Status: o.Status,
			Items: items, CreatedAt: o.CreatedAt,
		})
	}

	return &types.ListOrderResp{Orders: orders, Total: res.Total}, nil
}

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderReq) (resp *types.CancelOrderResp, err error) {
	userId := getUserIdFromCtx(l.ctx)
	res, err := l.svcCtx.OrderRpc.CancelOrder(l.ctx, &order2.CancelOrderRequest{
		OrderNo: req.OrderNo, UserId: userId,
	})
	if err != nil {
		return &types.CancelOrderResp{Success: false}, nil
	}
	return &types.CancelOrderResp{Success: res.Success}, nil
}

type PayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PayOrderLogic) PayOrder(req *types.PayOrderReq) (resp *types.PayOrderResp, err error) {
	userId := getUserIdFromCtx(l.ctx)
	res, err := l.svcCtx.OrderRpc.PayOrder(l.ctx, &order2.PayOrderRequest{
		OrderNo: req.OrderNo, UserId: userId, PayType: req.PayType,
	})
	if err != nil {
		return &types.PayOrderResp{Success: false}, nil
	}
	return &types.PayOrderResp{Success: res.Success}, nil
}
