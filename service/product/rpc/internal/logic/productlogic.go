package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"go-zero-ecommerce/common/errx"
	"go-zero-ecommerce/service/product/rpc/internal/svc"
	"go-zero-ecommerce/service/product/rpc/model"
	"go-zero-ecommerce/service/product/rpc/product"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductLogic) CreateProduct(in *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	if in.Name == "" || in.Price <= 0 {
		return nil, errx.ErrInvalidParam
	}

	p := &model.Product{
		CategoryId:    in.CategoryId,
		Name:          in.Name,
		Subtitle:      in.Subtitle,
		MainImage:     in.MainImage,
		SubImages:     in.SubImages,
		Detail:        in.Detail,
		Spec:          in.Spec,
		Price:         in.Price,
		OriginalPrice: in.OriginalPrice,
		Stock:         in.Stock,
		Brand:         in.Brand,
		Status:        1,
	}
	if p.OriginalPrice == 0 {
		p.OriginalPrice = p.Price
	}

	err := l.svcCtx.ProductModel.Insert(l.ctx, p)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	idStr := strconv.FormatInt(p.Id, 10)
	_ = l.svcCtx.BloomFilter.Add(l.ctx, []byte(idStr))
	l.indexToES(p)

	return &product.CreateProductResponse{
		Product: toProductInfo(p),
	}, nil
}

func (l *CreateProductLogic) indexToES(p *model.Product) {
	if l.svcCtx.ESClient == nil {
		return
	}
	_, err := l.svcCtx.ESClient.Index().
		Index("products").
		Id(strconv.FormatInt(p.Id, 10)).
		BodyJson(p).
		Do(l.ctx)
	if err != nil {
		logx.Errorf("ES index error: %v", err)
	}
}

type GetProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductLogic) GetProduct(in *product.GetProductRequest) (*product.GetProductResponse, error) {
	idStr := strconv.FormatInt(in.Id, 10)

	exists, err := l.svcCtx.BloomFilter.Contains(l.ctx, []byte(idStr))
	if err != nil {
		logx.Errorf("BloomFilter check error: %v", err)
	}
	if err == nil && !exists {
		return nil, errx.ErrProductNotFound
	}

	var p model.Product
	err = l.svcCtx.ProductCache.LoadOrStore(l.ctx, in.Id, &p, func() (interface{}, error) {
		dbProduct, dbErr := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
		if dbErr != nil {
			if errors.Is(dbErr, gorm.ErrRecordNotFound) {
				return nil, errx.ErrProductNotFound
			}
			return nil, dbErr
		}
		return dbProduct, nil
	})
	if err != nil {
		return nil, err
	}

	return &product.GetProductResponse{
		Product: toProductInfo(&p),
	}, nil
}

type UpdateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductLogic) UpdateProduct(in *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrProductNotFound
		}
		return nil, errx.ErrInternalServer
	}

	if in.CategoryId > 0 {
		p.CategoryId = in.CategoryId
	}
	if in.Name != "" {
		p.Name = in.Name
	}
	if in.Subtitle != "" {
		p.Subtitle = in.Subtitle
	}
	if in.MainImage != "" {
		p.MainImage = in.MainImage
	}
	if in.SubImages != "" {
		p.SubImages = in.SubImages
	}
	if in.Detail != "" {
		p.Detail = in.Detail
	}
	if in.Spec != "" {
		p.Spec = in.Spec
	}
	if in.Price > 0 {
		p.Price = in.Price
	}
	if in.OriginalPrice > 0 {
		p.OriginalPrice = in.OriginalPrice
	}
	if in.Status > 0 {
		p.Status = in.Status
	}
	if in.Brand != "" {
		p.Brand = in.Brand
	}

	err = l.svcCtx.ProductModel.Update(l.ctx, p)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	_ = l.svcCtx.ProductCache.Del(l.ctx, p.Id)
	l.updateES(p)

	return &product.UpdateProductResponse{
		Product: toProductInfo(p),
	}, nil
}

func (l *UpdateProductLogic) updateES(p *model.Product) {
	if l.svcCtx.ESClient == nil {
		return
	}
	_, err := l.svcCtx.ESClient.Update().
		Index("products").
		Id(strconv.FormatInt(p.Id, 10)).
		Doc(p).
		Do(l.ctx)
	if err != nil {
		logx.Errorf("ES update error: %v", err)
	}
}

type ListProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListProductLogic) ListProduct(in *product.ListProductRequest) (*product.ListProductResponse, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	var products []*model.Product
	var total int64
	var err error

	if in.Keyword != "" && l.svcCtx.ESClient != nil {
		products, total, err = l.searchByES(in)
	} else {
		products, total, err = l.svcCtx.ProductModel.List(l.ctx, in.Page, in.PageSize, in.CategoryId, in.Keyword)
	}
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	result := make([]*product.ProductInfo, 0, len(products))
	for _, p := range products {
		result = append(result, toProductInfo(p))
	}

	return &product.ListProductResponse{
		Products: result,
		Total:    total,
	}, nil
}

func (l *ListProductLogic) searchByES(in *product.ListProductRequest) ([]*model.Product, int64, error) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewTermQuery("status", 1))

	if in.Keyword != "" {
		multiMatch := elastic.NewMultiMatchQuery(in.Keyword, "name^3", "brand^2", "subtitle", "detail")
		boolQuery.Must(multiMatch)
	}
	if in.CategoryId > 0 {
		boolQuery.Filter(elastic.NewTermQuery("categoryId", in.CategoryId))
	}

	searchResult, err := l.svcCtx.ESClient.Search().
		Index("products").
		Query(boolQuery).
		From(int((in.Page - 1) * in.PageSize)).
		Size(int(in.PageSize)).
		Sort("_score", false).
		Pretty(true).
		Do(l.ctx)
	if err != nil {
		return nil, 0, err
	}

	total := searchResult.TotalHits()
	products := make([]*model.Product, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var p model.Product
		if err := json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, &p)
		}
	}
	_ = json.Marshal(boolQuery)
	fmt.Println("ES search total:", total)
	return products, total, nil
}

type GetProductListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductListLogic {
	return &GetProductListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductListLogic) GetProductList(in *product.GetProductListRequest) (*product.GetProductListResponse, error) {
	if len(in.Ids) == 0 {
		return &product.GetProductListResponse{Products: []*product.ProductInfo{}}, nil
	}

	products, err := l.svcCtx.ProductModel.ListByIds(l.ctx, in.Ids)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	result := make([]*product.ProductInfo, 0, len(products))
	for _, p := range products {
		result = append(result, toProductInfo(p))
	}

	return &product.GetProductListResponse{
		Products: result,
	}, nil
}

type ListCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoryLogic {
	return &ListCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCategoryLogic) ListCategory(in *product.ListCategoryRequest) (*product.ListCategoryResponse, error) {
	categories, err := l.svcCtx.CategoryModel.List(l.ctx, in.ParentId)
	if err != nil {
		return nil, errx.ErrInternalServer
	}

	result := make([]*product.CategoryInfo, 0, len(categories))
	for _, c := range categories {
		result = append(result, &product.CategoryInfo{
			Id:        c.Id,
			ParentId:  c.ParentId,
			Name:      c.Name,
			SortOrder: c.SortOrder,
			Status:    c.Status,
		})
	}

	return &product.ListCategoryResponse{
		Categories: result,
	}, nil
}

func toProductInfo(p *model.Product) *product.ProductInfo {
	return &product.ProductInfo{
		Id:            p.Id,
		CategoryId:    p.CategoryId,
		Name:          p.Name,
		Subtitle:      p.Subtitle,
		MainImage:     p.MainImage,
		SubImages:     p.SubImages,
		Detail:        p.Detail,
		Spec:          p.Spec,
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		Stock:         p.Stock,
		Sales:         p.Sales,
		Status:        p.Status,
		Brand:         p.Brand,
		CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
