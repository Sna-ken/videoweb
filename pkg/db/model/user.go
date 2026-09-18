package model

import "time"

type User struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id"`
	Nickname  string    `gorm:"column:nickname"`
	Avatar    string    `gorm:"column:avatar"`
	Signature string    `gorm:"column:signature"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
