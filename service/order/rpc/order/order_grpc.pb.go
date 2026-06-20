package order

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type OrderClient interface {
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	GetOrder(ctx context.Context, in *GetOrderRequest, opts ...grpc.CallOption) (*GetOrderResponse, error)
	ListOrder(ctx context.Context, in *ListOrderRequest, opts ...grpc.CallOption) (*ListOrderResponse, error)
	CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*CancelOrderResponse, error)
	PayOrder(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error)
}

type orderClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderClient(cc grpc.ClientConnInterface) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/order.Order/CreateOrder", in, out, opts...)
	return out, err
}

func (c *orderClient) GetOrder(ctx context.Context, in *GetOrderRequest, opts ...grpc.CallOption) (*GetOrderResponse, error) {
	out := new(GetOrderResponse)
	err := c.cc.Invoke(ctx, "/order.Order/GetOrder", in, out, opts...)
	return out, err
}

func (c *orderClient) ListOrder(ctx context.Context, in *ListOrderRequest, opts ...grpc.CallOption) (*ListOrderResponse, error) {
	out := new(ListOrderResponse)
	err := c.cc.Invoke(ctx, "/order.Order/ListOrder", in, out, opts...)
	return out, err
}

func (c *orderClient) CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*CancelOrderResponse, error) {
	out := new(CancelOrderResponse)
	err := c.cc.Invoke(ctx, "/order.Order/CancelOrder", in, out, opts...)
	return out, err
}

func (c *orderClient) PayOrder(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error) {
	out := new(PayOrderResponse)
	err := c.cc.Invoke(ctx, "/order.Order/PayOrder", in, out, opts...)
	return out, err
}

type OrderServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrder(context.Context, *GetOrderRequest) (*GetOrderResponse, error)
	ListOrder(context.Context, *ListOrderRequest) (*ListOrderResponse, error)
	CancelOrder(context.Context, *CancelOrderRequest) (*CancelOrderResponse, error)
	PayOrder(context.Context, *PayOrderRequest) (*PayOrderResponse, error)
	mustEmbedUnimplementedOrderServer()
}

type UnimplementedOrderServer struct{}

func (UnimplementedOrderServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedOrderServer) GetOrder(context.Context, *GetOrderRequest) (*GetOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedOrderServer) ListOrder(context.Context, *ListOrderRequest) (*ListOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedOrderServer) CancelOrder(context.Context, *CancelOrderRequest) (*CancelOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedOrderServer) PayOrder(context.Context, *PayOrderRequest) (*PayOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedOrderServer) mustEmbedUnimplementedOrderServer() {}

func RegisterOrderServer(s grpc.ServiceRegistrar, srv OrderServer) {
	s.RegisterService(&Order_ServiceDesc, srv)
}

var Order_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.Order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateOrder"},
		{MethodName: "GetOrder"},
		{MethodName: "ListOrder"},
		{MethodName: "CancelOrder"},
		{MethodName: "PayOrder"},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
