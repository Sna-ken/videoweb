package repository

import (
	"context"
	"errors"

	"github.com/Sna-ken/videoweb/pkg/db/model"
	"gorm.io/gorm"
)

func (r *AuthDB) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var account model.AuthAccount

	err := r.db.WithContext(ctx).
		Select("user_id").
		Where("username = ?", username).
		Take(&account).Error

	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}
