package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo      string    `gorm:"column:order_no;type:varchar(64);uniqueIndex;not null" json:"orderNo"`
	UserId       int64     `gorm:"column:user_id;index;not null" json:"userId"`
	TotalAmount  int64     `gorm:"column:total_amount;not null" json:"totalAmount"`
	PayAmount    int64     `gorm:"column:pay_amount;not null" json:"payAmount"`
	FreightAmount int64    `gorm:"column:freight_amount;default:0" json:"freightAmount"`
	DiscountAmount int64   `gorm:"column:discount_amount;default:0" json:"discountAmount"`
	Status       int32     `gorm:"column:status;default:0;index" json:"status"`
	PayType      int32     `gorm:"column:pay_type;default:0" json:"payType"`
	PayTime      *time.Time `gorm:"column:pay_time" json:"payTime"`
	ReceiverName    string `gorm:"column:receiver_name;type:varchar(64)" json:"receiverName"`
	ReceiverPhone   string `gorm:"column:receiver_phone;type:varchar(20)" json:"receiverPhone"`
	ReceiverAddress string `gorm:"column:receiver_address;type:varchar(500)" json:"receiverAddress"`
	Remark      string    `gorm:"column:remark;type:varchar(500)" json:"remark"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Order) TableName() string {
	return "order"
}

type OrderItem struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderId     int64     `gorm:"column:order_id;index;not null" json:"orderId"`
	OrderNo     string    `gorm:"column:order_no;type:varchar(64);index" json:"orderNo"`
	ProductId   int64     `gorm:"column:product_id;index;not null" json:"productId"`
	ProductName string    `gorm:"column:product_name;type:varchar(255);not null" json:"productName"`
	ProductImage string   `gorm:"column:product_image;type:varchar(500)" json:"productImage"`
	Price       int64     `gorm:"column:price;not null" json:"price"`
	Num         int32     `gorm:"column:num;not null" json:"num"`
	TotalPrice  int64     `gorm:"column:total_price;not null" json:"totalPrice"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (OrderItem) TableName() string {
	return "order_item"
}

type OrderModel interface {
	Insert(ctx context.Context, order *Order, items []*OrderItem) error
	FindByOrderNo(ctx context.Context, orderNo string) (*Order, error)
	FindById(ctx context.Context, id int64) (*Order, error)
	ListByUserId(ctx context.Context, userId int64, page, pageSize int32) ([]*Order, int64, error)
	UpdateStatus(ctx context.Context, orderNo string, oldStatus, newStatus int32) (int64, error)
}

type OrderItemModel interface {
	FindByOrderId(ctx context.Context, orderId int64) ([]*OrderItem, error)
}

type orderModel struct {
	db *gorm.DB
}

func NewOrderModel(db *gorm.DB) OrderModel {
	return &orderModel{db: db}
}

func (m *orderModel) Insert(ctx context.Context, order *Order, items []*OrderItem) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderId = order.Id
			items[i].OrderNo = order.OrderNo
		}
		if len(items) > 0 {
			if err := tx.Create(items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *orderModel) FindByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	var o Order
	err := m.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&o).Error
	return &o, err
}

func (m *orderModel) FindById(ctx context.Context, id int64) (*Order, error) {
	var o Order
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&o).Error
	return &o, err
}

func (m *orderModel) ListByUserId(ctx context.Context, userId int64, page, pageSize int32) ([]*Order, int64, error) {
	var orders []*Order
	var total int64
	err := m.db.WithContext(ctx).Model(&Order{}).Where("user_id = ?", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err = m.db.WithContext(ctx).Where("user_id = ?", userId).Order("id DESC").
		Offset(int(offset)).Limit(int(pageSize)).Find(&orders).Error
	return orders, total, err
}

func (m *orderModel) UpdateStatus(ctx context.Context, orderNo string, oldStatus, newStatus int32) (int64, error) {
	result := m.db.WithContext(ctx).Model(&Order{}).
		Where("order_no = ? AND status = ?", orderNo, oldStatus).
		Update("status", newStatus)
	return result.RowsAffected, result.Error
}

type orderItemModel struct {
	db *gorm.DB
}

func NewOrderItemModel(db *gorm.DB) OrderItemModel {
	return &orderItemModel{db: db}
}

func (m *orderItemModel) FindByOrderId(ctx context.Context, orderId int64) ([]*OrderItem, error) {
	var items []*OrderItem
	err := m.db.WithContext(ctx).Where("order_id = ?", orderId).Find(&items).Error
	return items, err
}
