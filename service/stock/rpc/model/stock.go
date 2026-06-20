package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductId int64     `gorm:"column:product_id;uniqueIndex;not null" json:"productId"`
	Total     int32     `gorm:"column:total;default:0" json:"total"`
	Available int32     `gorm:"column:available;default:0" json:"available"`
	LockStock int32     `gorm:"column:lock_stock;default:0" json:"lockStock"`
	Sales     int32     `gorm:"column:sales;default:0" json:"sales"`
	Version   int32     `gorm:"column:version;default:0" json:"version"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Stock) TableName() string {
	return "stock"
}

type StockModel interface {
	Insert(ctx context.Context, stock *Stock) error
	FindByProductId(ctx context.Context, productId int64) (*Stock, error)
	Update(ctx context.Context, stock *Stock) error
	DecrAvailable(ctx context.Context, productId int64, num int32, version int32) (int64, error)
	IncrAvailable(ctx context.Context, productId int64, num int32) (int64, error)
}

type stockModel struct {
	db *gorm.DB
}

func NewStockModel(db *gorm.DB) StockModel {
	return &stockModel{db: db}
}

func (m *stockModel) Insert(ctx context.Context, stock *Stock) error {
	return m.db.WithContext(ctx).Create(stock).Error
}

func (m *stockModel) FindByProductId(ctx context.Context, productId int64) (*Stock, error) {
	var stock Stock
	err := m.db.WithContext(ctx).Where("product_id = ?", productId).First(&stock).Error
	return &stock, err
}

func (m *stockModel) Update(ctx context.Context, stock *Stock) error {
	return m.db.WithContext(ctx).Save(stock).Error
}

func (m *stockModel) DecrAvailable(ctx context.Context, productId int64, num int32, version int32) (int64, error) {
	result := m.db.WithContext(ctx).Model(&Stock{}).
		Where("product_id = ? AND version = ? AND available >= ?", productId, version, num).
		Updates(map[string]interface{}{
			"available":  gorm.Expr("available - ?", num),
			"lock_stock": gorm.Expr("lock_stock + ?", num),
			"version":    gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}

func (m *stockModel) IncrAvailable(ctx context.Context, productId int64, num int32) (int64, error) {
	result := m.db.WithContext(ctx).Model(&Stock{}).
		Where("product_id = ?", productId).
		Updates(map[string]interface{}{
			"available":  gorm.Expr("available + ?", num),
			"lock_stock": gorm.Expr("lock_stock - ?", num),
		})
	return result.RowsAffected, result.Error
}
