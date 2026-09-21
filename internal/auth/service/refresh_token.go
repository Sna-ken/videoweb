package service

import (
	"context"
	"errors"

	"github.com/Sna-ken/videoweb/pkg/constants"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/redis/go-redis/v9"
)

func (s *AuthService) RefreshToken(ctx context.Context, userID, refreshToken string) (string, error) {
	claims, err := jwt.ValidateToken(refreshToken, constants.TypeRefreshToken)
	if err != nil {
		return "", err
	}
	if claims.Subject != userID {
		return "", errno.Wrap(errno.NewErr(errno.AuthInvalidErrorCode, "token用户不匹配"), nil)
	}

	_, err = s.authdb.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errno.Wrap(
				errno.AuthInvalidError,
				err,
			)
		}

		return "", errno.Wrap(
			errno.InternalDatabaseError,
			err,
		)
	}

	accesstoken, err := jwt.GenerateToken(userID, constants.TypeAccessToken)
	if err != nil {
		return "", err
	}

	return accesstoken, nil
}
