package model

import (
	"context"
	"time"
)

type Category struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentId  int64     `gorm:"column:parent_id;default:0;index" json:"parentId"`
	Name      string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	SortOrder int32     `gorm:"column:sort_order;default:0" json:"sortOrder"`
	Status    int32     `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Category) TableName() string {
	return "category"
}

type CategoryModel interface {
	Insert(ctx context.Context, category *Category) error
	FindOne(ctx context.Context, id int64) (*Category, error)
	List(ctx context.Context, parentId int64) ([]*Category, error)
}

type categoryModel struct {
	db *gorm.DB
}

func NewCategoryModel(db *gorm.DB) CategoryModel {
	return &categoryModel{db: db}
}

func (m *categoryModel) Insert(ctx context.Context, category *Category) error {
	return m.db.WithContext(ctx).Create(category).Error
}

func (m *categoryModel) FindOne(ctx context.Context, id int64) (*Category, error) {
	var category Category
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	return &category, err
}

func (m *categoryModel) List(ctx context.Context, parentId int64) ([]*Category, error) {
	var categories []*Category
	err := m.db.WithContext(ctx).Where("parent_id = ? AND status = ?", parentId, 1).
		Order("sort_order ASC, id ASC").Find(&categories).Error
	return categories, err
}
