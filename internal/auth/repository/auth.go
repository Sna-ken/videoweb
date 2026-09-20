package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthDB struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewAuthDB(db *gorm.DB, cache *redis.Client) *AuthDB {
	return &AuthDB{db: db, cache: cache}
}
