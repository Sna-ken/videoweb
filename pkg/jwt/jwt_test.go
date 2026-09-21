package jwt

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/pkg/constants"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/bytedance/mockey"
	gojwt "github.com/golang-jwt/jwt/v5"
)

func TestValidateExpiredAccessToken(t *testing.T) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate signing key: %v", err)
	}
	_, otherPrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate another signing key: %v", err)
	}

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal signing key: %v", err)
	}
	secret := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})

	now := time.Now()
	newToken := func(tokenType string, signingKey ed25519.PrivateKey) string {
		t.Helper()
		claims := Claims{
			Type: tokenType,
			RegisteredClaims: gojwt.RegisteredClaims{
				Subject:   "user-id",
				Issuer:    constants.Issuer,
				IssuedAt:  gojwt.NewNumericDate(now.Add(-2 * time.Minute)),
				ExpiresAt: gojwt.NewNumericDate(now.Add(-time.Minute)),
			},
		}
		token, signErr := gojwt.NewWithClaims(gojwt.SigningMethodEdDSA, claims).SignedString(signingKey)
		if signErr != nil {
			t.Fatalf("sign token: %v", signErr)
		}
		return token
	}

	tests := []struct {
		name     string
		token    string
		wantCode int64
	}{
		{
			name:  "expired access token is accepted",
			token: newToken(constants.TypeAccessToken, privateKey),
		},
		{
			name:     "refresh token cannot replace access token",
			token:    newToken(constants.TypeRefreshToken, privateKey),
			wantCode: errno.AuthInvalidErrorCode,
		},
		{
			name:     "invalid signature is rejected",
			token:    newToken(constants.TypeAccessToken, otherPrivateKey),
			wantCode: errno.AuthInvalidErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.PatchRun(func() {
				mockey.Mock(config.GetAuthorization).
					Return(&config.AuthorizationConfig{JWTSecret: string(secret)}, nil).
					Build()

				claims, validateErr := ValidateExpiredAccessToken(tt.token)
				if tt.wantCode != 0 {
					if validateErr == nil {
						t.Fatalf("ValidateExpiredAccessToken() error = nil, want code %d", tt.wantCode)
					}
					if got := errno.Convert(validateErr).Code(); got != tt.wantCode {
						t.Fatalf("ValidateExpiredAccessToken() code = %d, want %d", got, tt.wantCode)
					}
					return
				}

				if validateErr != nil {
					t.Fatalf("ValidateExpiredAccessToken() error = %v, want nil", validateErr)
				}
				if claims.Subject != "user-id" {
					t.Fatalf("ValidateExpiredAccessToken() subject = %q, want user-id", claims.Subject)
				}
			})
		})
	}
}
