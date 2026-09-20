package jwt

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/pkg/constans"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, string, error) {
	cfg, err := config.GetAuthorization()
	if err != nil || cfg == nil {
		return "", "", errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "获取配置失败"), err)
	}
	if userID == "" {
		return "", "", errno.Wrap(errno.NewErr(errno.ParamEmptyErrorCode, "用户id为空"), nil)
	}

	accessToken, err := createToken(userID, "access_token", cfg)
	if err != nil {
		return "", "", errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "生成access_token失败"), err)
	}

	refreshToken, err := createToken(userID, "refresh_token", cfg)
	if err != nil {
		return "", "", errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "生成refresh_token失败"), err)
	}

	return accessToken, refreshToken, nil
}

func createToken(userID string, tokentype string, cfg *config.AuthorizationConfig) (string, error) {
	now := time.Now()
	var exp time.Time

	switch tokentype {
	case "access_token":
		exp = now.Add(time.Duration(constans.AccessTokenExpire))
	case "refresh_token":
		exp = now.Add(time.Duration(constans.RefreshTokenExpire))
	default:
		return "", errno.NewErr(errno.ParamInvalidErrorCode, "token类型无效")
	}

	claims := Claims{
		Type: tokentype,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        uuid.NewString(),
			Issuer:    constans.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	key, err := jwt.ParseEdPrivateKeyFromPEM([]byte(cfg.JWTSecret))
	if err != nil {
		return "", errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "解析私钥失败"), err)
	}
	return unsignedToken.SignedString(key)
}

func ValidateToken(token string, tokenType string) (*Claims, error) {
	if token == "" {
		return nil, errno.Wrap(errno.NewErr(errno.ParamEmptyErrorCode, "token为空"), nil)
	}
	if tokenType != constans.TypeAccessToken && tokenType != constans.TypeRefreshToken {
		return nil, errno.Wrap(errno.NewErr(errno.ParamInvalidErrorCode, "token类型无效"), nil)
	}

	cfg, err := config.GetAuthorization()
	if err != nil || cfg == nil {
		return nil, errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "获取配置失败"), err)
	}
	publicKey, err := getPublicKey(cfg.JWTSecret)
	if err != nil {
		return nil, errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "获取公钥失败"), err)
	}

	parsedToken, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(_ *jwt.Token) (any, error) {
			return publicKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithIssuer(constans.Issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errno.Wrap(errno.NewErr(errno.AuthExpiredErrorCode, "token已过期"), err)
		}
		return nil, errno.Wrap(errno.NewErr(errno.AuthInvalidErrorCode, "token无效"), err)
	}
	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, errno.Wrap(errno.NewErr(errno.AuthInvalidErrorCode, "token无效"), nil)
	}
	if claims.Type != tokenType {
		return nil, errno.Wrap(errno.NewErr(errno.AuthInvalidErrorCode, "token类型不匹配"), nil)
	}

	return claims, nil
}

func getPublicKey(secret string) (ed25519.PublicKey, error) {
	key, err := jwt.ParseEdPrivateKeyFromPEM([]byte(secret))
	if err != nil {
		return nil, errno.Wrap(errno.InternalServiceError, err)
	}

	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errno.Wrap(
			errno.InternalServiceError,
			fmt.Errorf("JWT私钥类型错误: %T", key),
		)
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)

	return publicKey, nil
}
