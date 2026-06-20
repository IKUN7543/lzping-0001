package model

import (
	"time"
)

type User struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(64);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Nickname  string    `gorm:"column:nickname;type:varchar(64)" json:"nickname"`
	Mobile    string    `gorm:"column:mobile;type:varchar(20);uniqueIndex" json:"mobile"`
	Email     string    `gorm:"column:email;type:varchar(128);uniqueIndex" json:"email"`
	Gender    int32     `gorm:"column:gender;default:0" json:"gender"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	Status    int32     `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}
