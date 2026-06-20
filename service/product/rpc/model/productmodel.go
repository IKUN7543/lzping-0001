package model

import (
	"context"
	"gorm.io/gorm"
)

type ProductModel interface {
	Insert(ctx context.Context, product *Product) error
	FindOne(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, pageSize int32, categoryId int64, keyword string) ([]*Product, int64, error)
	ListByIds(ctx context.Context, ids []int64) ([]*Product, error)
}

type productModel struct {
	db *gorm.DB
}

func NewProductModel(db *gorm.DB) ProductModel {
	return &productModel{db: db}
}

func (m *productModel) Insert(ctx context.Context, product *Product) error {
	return m.db.WithContext(ctx).Create(product).Error
}

func (m *productModel) FindOne(ctx context.Context, id int64) (*Product, error) {
	var product Product
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *productModel) Update(ctx context.Context, product *Product) error {
	return m.db.WithContext(ctx).Save(product).Error
}

func (m *productModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Delete(&Product{}).Error
}

func (m *productModel) List(ctx context.Context, page, pageSize int32, categoryId int64, keyword string) ([]*Product, int64, error) {
	var products []*Product
	var total int64
	query := m.db.WithContext(ctx).Model(&Product{})

	if categoryId > 0 {
		query = query.Where("category_id = ?", categoryId)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR brand LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query = query.Where("status = ?", 1)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (m *productModel) ListByIds(ctx context.Context, ids []int64) ([]*Product, error) {
	var products []*Product
	err := m.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error
	return products, err
}
