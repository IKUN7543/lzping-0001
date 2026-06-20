package stock

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type StockClient interface {
	CreateStock(ctx context.Context, in *CreateStockRequest, opts ...grpc.CallOption) (*CreateStockResponse, error)
	GetStock(ctx context.Context, in *GetStockRequest, opts ...grpc.CallOption) (*GetStockResponse, error)
	DeductStock(ctx context.Context, in *DeductStockRequest, opts ...grpc.CallOption) (*DeductStockResponse, error)
	ReturnStock(ctx context.Context, in *ReturnStockRequest, opts ...grpc.CallOption) (*ReturnStockResponse, error)
	ConfirmStock(ctx context.Context, in *ConfirmStockRequest, opts ...grpc.CallOption) (*ConfirmStockResponse, error)
}

type stockClient struct {
	cc grpc.ClientConnInterface
}

func NewStockClient(cc grpc.ClientConnInterface) StockClient {
	return &stockClient{cc}
}

func (c *stockClient) CreateStock(ctx context.Context, in *CreateStockRequest, opts ...grpc.CallOption) (*CreateStockResponse, error) {
	out := new(CreateStockResponse)
	err := c.cc.Invoke(ctx, "/stock.Stock/CreateStock", in, out, opts...)
	return out, err
}

func (c *stockClient) GetStock(ctx context.Context, in *GetStockRequest, opts ...grpc.CallOption) (*GetStockResponse, error) {
	out := new(GetStockResponse)
	err := c.cc.Invoke(ctx, "/stock.Stock/GetStock", in, out, opts...)
	return out, err
}

func (c *stockClient) DeductStock(ctx context.Context, in *DeductStockRequest, opts ...grpc.CallOption) (*DeductStockResponse, error) {
	out := new(DeductStockResponse)
	err := c.cc.Invoke(ctx, "/stock.Stock/DeductStock", in, out, opts...)
	return out, err
}

func (c *stockClient) ReturnStock(ctx context.Context, in *ReturnStockRequest, opts ...grpc.CallOption) (*ReturnStockResponse, error) {
	out := new(ReturnStockResponse)
	err := c.cc.Invoke(ctx, "/stock.Stock/ReturnStock", in, out, opts...)
	return out, err
}

func (c *stockClient) ConfirmStock(ctx context.Context, in *ConfirmStockRequest, opts ...grpc.CallOption) (*ConfirmStockResponse, error) {
	out := new(ConfirmStockResponse)
	err := c.cc.Invoke(ctx, "/stock.Stock/ConfirmStock", in, out, opts...)
	return out, err
}

type StockServer interface {
	CreateStock(context.Context, *CreateStockRequest) (*CreateStockResponse, error)
	GetStock(context.Context, *GetStockRequest) (*GetStockResponse, error)
	DeductStock(context.Context, *DeductStockRequest) (*DeductStockResponse, error)
	ReturnStock(context.Context, *ReturnStockRequest) (*ReturnStockResponse, error)
	ConfirmStock(context.Context, *ConfirmStockRequest) (*ConfirmStockResponse, error)
	mustEmbedUnimplementedStockServer()
}

type UnimplementedStockServer struct{}

func (UnimplementedStockServer) CreateStock(context.Context, *CreateStockRequest) (*CreateStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedStockServer) GetStock(context.Context, *GetStockRequest) (*GetStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedStockServer) DeductStock(context.Context, *DeductStockRequest) (*DeductStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedStockServer) ReturnStock(context.Context, *ReturnStockRequest) (*ReturnStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedStockServer) ConfirmStock(context.Context, *ConfirmStockRequest) (*ConfirmStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
func (UnimplementedStockServer) mustEmbedUnimplementedStockServer() {}

func RegisterStockServer(s grpc.ServiceRegistrar, srv StockServer) {
	s.RegisterService(&Stock_ServiceDesc, srv)
}

var Stock_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stock.Stock",
	HandlerType: (*StockServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateStock"},
		{MethodName: "GetStock"},
		{MethodName: "DeductStock"},
		{MethodName: "ReturnStock"},
		{MethodName: "ConfirmStock"},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stock.proto",
}
