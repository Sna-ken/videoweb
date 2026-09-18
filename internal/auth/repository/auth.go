package repository

import "gorm.io/gorm"

type AuthDB struct {
	db *gorm.DB
}

func NewAuthDB(db *gorm.DB) *AuthDB {
	return &AuthDB{db: db}
}
