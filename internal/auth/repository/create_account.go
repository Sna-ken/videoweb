package repository

import (
	"context"
	"errors"

	"github.com/Sna-ken/videoweb/pkg/db/model"
	"gorm.io/gorm"
)

var ErrAccountAlreadyExists = errors.New("auth account already exists")

func (r *AuthDB) CreateAccount(ctx context.Context, account *model.AuthAccount) error {
	err := r.db.WithContext(ctx).Create(account).Error
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrAccountAlreadyExists
	}

	return err
}
