package model

import (
	"time"
)

type Product struct {
	Id            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CategoryId    int64     `gorm:"column:category_id;index;not null" json:"categoryId"`
	Name          string    `gorm:"column:name;type:varchar(255);index;not null" json:"name"`
	Subtitle      string    `gorm:"column:subtitle;type:varchar(500)" json:"subtitle"`
	MainImage     string    `gorm:"column:main_image;type:varchar(500)" json:"mainImage"`
	SubImages     string    `gorm:"column:sub_images;type:text" json:"subImages"`
	Detail        string    `gorm:"column:detail;type:text" json:"detail"`
	Spec          string    `gorm:"column:spec;type:varchar(500)" json:"spec"`
	Price         int64     `gorm:"column:price;not null" json:"price"`
	OriginalPrice int64     `gorm:"column:original_price;not null" json:"originalPrice"`
	Stock         int32     `gorm:"column:stock;default:0" json:"stock"`
	Sales         int32     `gorm:"column:sales;default:0" json:"sales"`
	Status        int32     `gorm:"column:status;default:1" json:"status"`
	Brand         string    `gorm:"column:brand;type:varchar(128);index" json:"brand"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Product) TableName() string {
	return "product"
}
