package logic

import (
	"context"
	"errors"
	"fmt"
	"go-zero-ecommerce/common/errx"
	"go-zero-ecommerce/common/kafkaq"
	"go-zero-ecommerce/service/order/rpc/internal/svc"
	"go-zero-ecommerce/service/order/rpc/model"
	"go-zero-ecommerce/service/order/rpc/order"
	product2 "go-zero-ecommerce/service/product/rpc/product"
	stock2 "go-zero-ecommerce/service/stock/rpc/stock"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if len(in.Items) == 0 {
		return nil, errx.ErrInvalidParam
	}

	orderId := l.svcCtx.SnowWorker.Id()
	orderNo := fmt.Sprintf("SO%s%010d", time.Now().Format("20060102150405"), orderId%10000000000)

	productIds := make([]int64, 0, len(in.Items))
	for _, it := range in.Items {
		productIds = append(productIds, it.ProductId)
	}
	products, err := l.svcCtx.ProductRpc.GetProductList(l.ctx, &product2.GetProductListRequest{Ids: productIds})
	if err != nil {
		logx.Errorf("Get product list error: %v", err)
	}

	var totalAmount int64
	msgItems := make([]svc.OrderItemKafkaMsg, 0, len(in.Items))
	for i, it := range in.Items {
		price := it.Price
		if price == 0 && products != nil && i < len(products.Products) {
			price = products.Products[i].Price
		}
		totalAmount += price * int64(it.Num)
		msgItems = append(msgItems, svc.OrderItemKafkaMsg{
			ProductId: it.ProductId, Num: it.Num, Price: price,
		})
	}
	payAmount := totalAmount

	for i, it := range in.Items {
		deductRes, err := l.svcCtx.StockRpc.DeductStock(l.ctx, &stock2.DeductStockRequest{
			ProductId: it.ProductId, Num: it.Num, OrderNo: orderNo,
		})
		if err != nil || !deductRes.Success {
			for j := 0; j < i; j++ {
				prev := in.Items[j]
				l.svcCtx.StockRpc.ReturnStock(l.ctx, &stock2.ReturnStockRequest{
					ProductId: prev.ProductId, Num: prev.Num, OrderNo: orderNo,
				})
			}
			return nil, errx.ErrStockNotEnough
		}
	}

	kafkaMsg := svc.OrderKafkaMsg{
		OrderNo:         orderNo,
		UserId:          in.UserId,
		Items:           msgItems,
		TotalAmount:     totalAmount,
		PayAmount:       payAmount,
		ReceiverName:    in.ReceiverName,
		ReceiverPhone:   in.ReceiverPhone,
		ReceiverAddress: in.ReceiverAddress,
		Remark:          in.Remark,
	}

	if l.svcCtx.KafkaProducer != nil {
		brokers := strings.Split(l.svcCtx.Config.Kafka.Brokers, ",")
		kafkaProducer := kafkaq.NewProducer(brokers)
		err := kafkaProducer.Send(l.ctx, l.svcCtx.Config.Kafka.Topic, orderNo, kafkaMsg)
		if err != nil {
			logx.Errorf("Kafka send error: %v, fallback to sync", err)
			l.createOrderSync(kafkaMsg, msgItems)
		}
		kafkaProducer.Close()
	} else {
		l.createOrderSync(kafkaMsg, msgItems)
	}

	return &order.CreateOrderResponse{OrderNo: orderNo, Id: orderId, PayAmount: payAmount}, nil
}

func (l *CreateOrderLogic) createOrderSync(kafkaMsg svc.OrderKafkaMsg, items []svc.OrderItemKafkaMsg) {
	o := &model.Order{
		OrderNo:         kafkaMsg.OrderNo,
		UserId:          kafkaMsg.UserId,
		TotalAmount:     kafkaMsg.TotalAmount,
		PayAmount:       kafkaMsg.PayAmount,
		Status:          1,
		ReceiverName:    kafkaMsg.ReceiverName,
		ReceiverPhone:   kafkaMsg.ReceiverPhone,
		ReceiverAddress: kafkaMsg.ReceiverAddress,
		Remark:          kafkaMsg.Remark,
	}
	orderItems := make([]*model.OrderItem, 0, len(items))
	for _, it := range items {
		orderItems = append(orderItems, &model.OrderItem{
			ProductId: it.ProductId, Num: it.Num, Price: it.Price,
			TotalPrice: it.Price * int64(it.Num),
		})
	}
	err := l.svcCtx.OrderModel.Insert(l.ctx, o, orderItems)
	if err != nil {
		logx.Errorf("Sync create order error: %v", err)
	}
}

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	o, err := l.svcCtx.OrderModel.FindByOrderNo(l.ctx, in.OrderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrOrderNotFound
		}
		return nil, errx.ErrInternalServer
	}

	items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, o.Id)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	orderItems := make([]*order.OrderItemInfo, 0, len(items))
	for _, it := range items {
		orderItems = append(orderItems, &order.OrderItemInfo{
			Id: it.Id, ProductId: it.ProductId, ProductName: it.ProductName,
			ProductImage: it.ProductImage, Price: it.Price, Num: it.Num, TotalPrice: it.TotalPrice,
		})
	}

	return &order.GetOrderResponse{Order: &order.OrderInfo{
		Id: o.Id, OrderNo: o.OrderNo, UserId: o.UserId, TotalAmount: o.TotalAmount,
		PayAmount: o.PayAmount, FreightAmount: o.FreightAmount, DiscountAmount: o.DiscountAmount,
		Status: o.Status, PayType: o.PayType, ReceiverName: o.ReceiverName,
		ReceiverPhone: o.ReceiverPhone, ReceiverAddress: o.ReceiverAddress,
		Remark: o.Remark, Items: orderItems,
		CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"), UpdatedAt: o.UpdatedAt.Format("2006-01-02 15:04:05"),
	}}, nil
}

type ListOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderLogic {
	return &ListOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListOrderLogic) ListOrder(in *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	orders, total, err := l.svcCtx.OrderModel.ListByUserId(l.ctx, in.UserId, in.Page, in.PageSize)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	result := make([]*order.OrderInfo, 0, len(orders))
	for _, o := range orders {
		items, _ := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, o.Id)
		orderItems := make([]*order.OrderItemInfo, 0, len(items))
		for _, it := range items {
			orderItems = append(orderItems, &order.OrderItemInfo{
				Id: it.Id, ProductId: it.ProductId, ProductName: it.ProductName, Num: it.Num,
				Price: it.Price, TotalPrice: it.TotalPrice, ProductImage: it.ProductImage,
			})
		}
		result = append(result, &order.OrderInfo{
			Id: o.Id, OrderNo: o.OrderNo, UserId: o.UserId, TotalAmount: o.TotalAmount,
			PayAmount: o.PayAmount, Status: o.Status, Items: orderItems,
			CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &order.ListOrderResponse{Orders: result, Total: total}, nil
}

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	o, err := l.svcCtx.OrderModel.FindByOrderNo(l.ctx, in.OrderNo)
	if err != nil {
		return &order.CancelOrderResponse{Success: false}, errx.ErrOrderNotFound
	}
	if o.UserId != in.UserId {
		return &order.CancelOrderResponse{Success: false}, errx.ErrUnauthorized
	}

	affected, err := l.svcCtx.OrderModel.UpdateStatus(l.ctx, in.OrderNo, 1, 4)
	if err != nil || affected == 0 {
		return &order.CancelOrderResponse{Success: false}, errx.ErrOrderCreateFail
	}

	items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, o.Id)
	if err == nil {
		for _, it := range items {
			l.svcCtx.StockRpc.ReturnStock(l.ctx, &stock2.ReturnStockRequest{
				ProductId: it.ProductId, Num: it.Num, OrderNo: in.OrderNo,
			})
		}
	}

	return &order.CancelOrderResponse{Success: true}, nil
}

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *PayOrderLogic) PayOrder(in *order.PayOrderRequest) (*order.PayOrderResponse, error) {
	o, err := l.svcCtx.OrderModel.FindByOrderNo(l.ctx, in.OrderNo)
	if err != nil {
		return &order.PayOrderResponse{Success: false}, errx.ErrOrderNotFound
	}
	if o.UserId != in.UserId {
		return &order.PayOrderResponse{Success: false}, errx.ErrUnauthorized
	}

	affected, err := l.svcCtx.OrderModel.UpdateStatus(l.ctx, in.OrderNo, 1, 2)
	if err != nil || affected == 0 {
		return &order.PayOrderResponse{Success: false}, errx.ErrOrderCreateFail
	}

	items, _ := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, o.Id)
	for _, it := range items {
		l.svcCtx.StockRpc.ConfirmStock(l.ctx, &stock2.ConfirmStockRequest{
			ProductId: it.ProductId, Num: it.Num, OrderNo: in.OrderNo,
		})
	}

	return &order.PayOrderResponse{Success: true}, nil
}
