package repository

import (
	"context"

	"github.com/Sna-ken/videoweb/pkg/db/model"
)

func (r *UserDB) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
