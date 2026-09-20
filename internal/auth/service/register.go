package service

import (
	"context"
	"errors"

	"github.com/Sna-ken/videoweb/internal/auth/repository"
	"github.com/Sna-ken/videoweb/pkg/db/model"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/Sna-ken/videoweb/pkg/utils"
	"github.com/google/uuid"
)

func (s *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errno.Wrap(errno.ParamEmptyError, nil)
	}

	if valid, err := utils.CheckPasswordValid(password); !valid {
		return "", errno.Wrap(errno.NewErr(errno.PasswordInvalidErrorCode, "密码无效"), err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", errno.Wrap(errno.NewErr(errno.PasswordHashErrorCode, "加密失败"), err)
	}

	exists, err := s.authdb.ExistsByUsername(ctx, username)
	if err != nil {
		return "", errno.Wrap(errno.InternalDatabaseError, err)
	}
	if exists {
		return "", errno.Wrap(
			errno.NewErr(errno.UserHasExistedErrorCode, "用户已存在"), nil)
	}

	userID := uuid.NewString()
	account := &model.AuthAccount{
		UserID:   userID,
		Username: username,
		Password: hashedPassword,
	}

	err = s.authdb.CreateAccount(ctx, account)
	switch {
	case err == nil:
		return userID, nil
	case errors.Is(err, repository.ErrAccountAlreadyExists):
		return "", errno.Wrap(
			errno.NewErr(errno.UserHasExistedErrorCode, "用户已存在"),
			err,
		)
	default:
		return "", errno.Wrap(errno.InternalDatabaseError, err)
	}
}
