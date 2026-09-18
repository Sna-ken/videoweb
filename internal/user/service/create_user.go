package service

import (
	"context"
	"errors"

	"github.com/Sna-ken/videoweb/pkg/db/model"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *UserService) CreateUser(ctx context.Context, userID, username string) error {
	if userID == "" || username == "" {
		return errno.Wrap(errno.ParamEmptyError, nil)
	}

	user := &model.User{
		ID:       uuid.NewString(),
		UserID:   userID,
		Nickname: username,
	}

	err := s.userdb.CreateUser(ctx, user)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return errno.Wrap(
			errno.NewErr(errno.UserHasExistedErrorCode, "用户已存在"),
			err,
		)
	default:
		return errno.Wrap(errno.InternalDatabaseError, err)
	}
}
