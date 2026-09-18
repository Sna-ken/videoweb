package repository

import (
	"context"

	"github.com/Sna-ken/videoweb/pkg/db/model"
)

func (r *AuthDB) SearchByName(ctx context.Context, username string) (*model.AuthAccount, error) {
	var account model.AuthAccount
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
