package repository

import "context"

func (r *AuthDB) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := GetRefreshTokenKey(refreshToken)
	return r.cache.Del(ctx, key).Err()
}
