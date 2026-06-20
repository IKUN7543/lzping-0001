package logic

import (
	"context"
	"go-zero-ecommerce/service/product/api/internal/svc"
	"go-zero-ecommerce/service/product/api/internal/types"
	"go-zero-ecommerce/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoLogic {
	return &ProductInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductInfoLogic) ProductInfo(req *types.ProductInfoReq) (resp *types.ProductInfoResp, err error) {
	res, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.ProductInfoResp{
		Id:            res.Product.Id,
		CategoryId:    res.Product.CategoryId,
		Name:          res.Product.Name,
		Subtitle:      res.Product.Subtitle,
		MainImage:     res.Product.MainImage,
		SubImages:     res.Product.SubImages,
		Detail:        res.Product.Detail,
		Spec:          res.Product.Spec,
		Price:         res.Product.Price,
		OriginalPrice: res.Product.OriginalPrice,
		Stock:         res.Product.Stock,
		Sales:         res.Product.Sales,
		Status:        res.Product.Status,
		Brand:         res.Product.Brand,
		CreatedAt:     res.Product.CreatedAt,
		UpdatedAt:     res.Product.UpdatedAt,
	}, nil
}

type ListProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductLogic) ListProduct(req *types.ListProductReq) (resp *types.ListProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.ListProduct(l.ctx, &product.ListProductRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		CategoryId: req.CategoryId,
		Keyword:    req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.ProductItem, 0, len(res.Products))
	for _, p := range res.Products {
		items = append(items, types.ProductItem{
			Id:            p.Id,
			CategoryId:    p.CategoryId,
			Name:          p.Name,
			Subtitle:      p.Subtitle,
			MainImage:     p.MainImage,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			Sales:         p.Sales,
			Brand:         p.Brand,
		})
	}

	return &types.ListProductResp{
		Products: items,
		Total:    res.Total,
	}, nil
}

type SearchProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductLogic {
	return &SearchProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchProductLogic) SearchProduct(req *types.SearchProductReq) (resp *types.SearchProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.ListProduct(l.ctx, &product.ListProductRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		CategoryId: req.CategoryId,
		Keyword:    req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.ProductItem, 0, len(res.Products))
	for _, p := range res.Products {
		items = append(items, types.ProductItem{
			Id:            p.Id,
			CategoryId:    p.CategoryId,
			Name:          p.Name,
			Subtitle:      p.Subtitle,
			MainImage:     p.MainImage,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			Sales:         p.Sales,
			Brand:         p.Brand,
		})
	}

	return &types.SearchProductResp{
		Products: items,
		Total:    res.Total,
	}, nil
}

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (resp *types.CreateProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.CreateProduct(l.ctx, &product.CreateProductRequest{
		CategoryId:    req.CategoryId,
		Name:          req.Name,
		Subtitle:      req.Subtitle,
		MainImage:     req.MainImage,
		SubImages:     req.SubImages,
		Detail:        req.Detail,
		Spec:          req.Spec,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Brand:         req.Brand,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateProductResp{
		Id:         res.Product.Id,
		CategoryId: res.Product.CategoryId,
		Name:       res.Product.Name,
		Price:      res.Product.Price,
	}, nil
}

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductReq) (resp *types.UpdateProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.UpdateProduct(l.ctx, &product.UpdateProductRequest{
		Id:            req.Id,
		CategoryId:    req.CategoryId,
		Name:          req.Name,
		Subtitle:      req.Subtitle,
		MainImage:     req.MainImage,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Status:        req.Status,
		Brand:         req.Brand,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateProductResp{
		Id:     res.Product.Id,
		Name:   res.Product.Name,
		Status: res.Product.Status,
	}, nil
}

type ListCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoryLogic {
	return &ListCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCategoryLogic) ListCategory(req *types.ListCategoryReq) (resp *types.ListCategoryResp, err error) {
	res, err := l.svcCtx.ProductRpc.ListCategory(l.ctx, &product.ListCategoryRequest{
		ParentId: req.ParentId,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.CategoryItem, 0, len(res.Categories))
	for _, c := range res.Categories {
		items = append(items, types.CategoryItem{
			Id:        c.Id,
			ParentId:  c.ParentId,
			Name:      c.Name,
			SortOrder: c.SortOrder,
		})
	}

	return &types.ListCategoryResp{
		Categories: items,
	}, nil
}
