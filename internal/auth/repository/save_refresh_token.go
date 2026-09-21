package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Sna-ken/videoweb/pkg/constants"
)

func (r *AuthDB) SaveRefreshToken(ctx context.Context, refreshToken string) error {
	key := GetRefreshTokenKey(refreshToken)
	return r.cache.Set(ctx, key, 1, constants.RefreshTokenExpire).Err()
}

func GetRefreshTokenKey(refreshToken string) string {
	sum := sha256.Sum256([]byte(refreshToken))
	return constants.RefreshTokenKeyPrefix + hex.EncodeToString(sum[:])
}
