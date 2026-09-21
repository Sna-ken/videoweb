package service

import (
	"context"
	"errors"

	"github.com/Sna-ken/videoweb/pkg/errno"
)

func (s *AuthService) Logout(ctx context.Context, refreshtoken string) error {
	if refreshtoken == "" {
		return errno.Wrap(errno.NewErr(errno.ParamEmptyErrorCode, "refresh token 为空"), errors.New("refresh token 为空"))
	}

	if err := s.authdb.DeleteRefreshToken(ctx, refreshtoken); err != nil {
		return errno.Wrap(errno.InternalDatabaseError, err)
	}

	return nil
}
