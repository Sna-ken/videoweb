package repository

import "context"

func (r *AuthDB) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := GetRefreshTokenKey(refreshToken)
	return r.cache.Get(ctx, key).Result()
}
