package server

import (
	"context"
	"go-zero-ecommerce/service/product/rpc/internal/logic"
	"go-zero-ecommerce/service/product/rpc/internal/svc"
	"go-zero-ecommerce/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductServer struct {
	svcCtx *svc.ServiceContext
	product.UnimplementedProductServer
}

func NewProductServer(svcCtx *svc.ServiceContext) *ProductServer {
	return &ProductServer{
		svcCtx: svcCtx,
	}
}

func (s *ProductServer) CreateProduct(ctx context.Context, in *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	l := logic.NewCreateProductLogic(ctx, s.svcCtx)
	logx.Infof("CreateProduct request: name=%s", in.Name)
	return l.CreateProduct(in)
}

func (s *ProductServer) GetProduct(ctx context.Context, in *product.GetProductRequest) (*product.GetProductResponse, error) {
	l := logic.NewGetProductLogic(ctx, s.svcCtx)
	return l.GetProduct(in)
}

func (s *ProductServer) UpdateProduct(ctx context.Context, in *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	l := logic.NewUpdateProductLogic(ctx, s.svcCtx)
	logx.Infof("UpdateProduct request: id=%d", in.Id)
	return l.UpdateProduct(in)
}

func (s *ProductServer) ListProduct(ctx context.Context, in *product.ListProductRequest) (*product.ListProductResponse, error) {
	l := logic.NewListProductLogic(ctx, s.svcCtx)
	return l.ListProduct(in)
}

func (s *ProductServer) GetProductList(ctx context.Context, in *product.GetProductListRequest) (*product.GetProductListResponse, error) {
	l := logic.NewGetProductListLogic(ctx, s.svcCtx)
	return l.GetProductList(in)
}

func (s *ProductServer) ListCategory(ctx context.Context, in *product.ListCategoryRequest) (*product.ListCategoryResponse, error) {
	l := logic.NewListCategoryLogic(ctx, s.svcCtx)
	return l.ListCategory(in)
}
