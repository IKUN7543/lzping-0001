package model

import (
	"context"
	"gorm.io/gorm"
)

type UserModel interface {
	Insert(ctx context.Context, user *User) error
	FindOne(ctx context.Context, id int64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByMobile(ctx context.Context, mobile string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type userModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) UserModel {
	return &userModel{db: db}
}

func (m *userModel) Insert(ctx context.Context, user *User) error {
	return m.db.WithContext(ctx).Create(user).Error
}

func (m *userModel) FindOne(ctx context.Context, id int64) (*User, error) {
	var user User
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *userModel) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := m.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *userModel) FindByMobile(ctx context.Context, mobile string) (*User, error) {
	var user User
	err := m.db.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *userModel) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *userModel) Update(ctx context.Context, user *User) error {
	return m.db.WithContext(ctx).Save(user).Error
}
