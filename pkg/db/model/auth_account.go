package model

import "time"

type AuthAccount struct {
	UserID     string    `gorm:"column:user_id;primaryKey"`
	Username   string    `gorm:"column:username"`
	Password   string    `gorm:"column:password"`
	MFAEnabled bool      `gorm:"column:mfa_enabled"`
	MFASecret  string    `gorm:"column:mfa_secret"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (AuthAccount) TableName() string {
	return "auth_accounts"
}
