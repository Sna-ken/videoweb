package repository

import (
	"context"

	"github.com/Sna-ken/videoweb/pkg/db/model"
)

func (r *UserDB) CountUser(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
