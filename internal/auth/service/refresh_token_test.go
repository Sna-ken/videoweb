package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sna-ken/videoweb/internal/auth/repository"
	"github.com/Sna-ken/videoweb/pkg/constants"
	"github.com/Sna-ken/videoweb/pkg/errno"
	jwtutil "github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/bytedance/mockey"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func TestRefreshToken(t *testing.T) {
	validateErr := errno.Wrap(errno.AuthInvalidError, errors.New("invalid refresh token"))
	databaseErr := errors.New("redis unavailable")
	generateErr := errno.Wrap(errno.InternalServiceError, errors.New("generate token failed"))

	tests := []struct {
		name              string
		requestUserID     string
		claimsUserID      string
		validateErr       error
		getErr            error
		generateErr       error
		wantAccessToken   string
		wantCode          int64
		wantCause         error
		wantGetCalls      int
		wantGenerateCalls int
	}{
		{
			name:          "token validation failed",
			requestUserID: "user-id",
			validateErr:   validateErr,
			wantCode:      errno.AuthInvalidErrorCode,
			wantCause:     validateErr,
		},
		{
			name:          "access and refresh token users do not match",
			requestUserID: "user-id",
			claimsUserID:  "other-user-id",
			wantCode:      errno.AuthInvalidErrorCode,
		},
		{
			name:          "refresh token session does not exist",
			requestUserID: "user-id",
			claimsUserID:  "user-id",
			getErr:        redis.Nil,
			wantCode:      errno.AuthInvalidErrorCode,
			wantCause:     redis.Nil,
			wantGetCalls:  1,
		},
		{
			name:          "refresh token query failed",
			requestUserID: "user-id",
			claimsUserID:  "user-id",
			getErr:        databaseErr,
			wantCode:      errno.InternalDatabaseErrorCode,
			wantCause:     databaseErr,
			wantGetCalls:  1,
		},
		{
			name:              "access token generation failed",
			requestUserID:     "user-id",
			claimsUserID:      "user-id",
			generateErr:       generateErr,
			wantCode:          errno.InternalServiceErrorCode,
			wantCause:         generateErr,
			wantGetCalls:      1,
			wantGenerateCalls: 1,
		},
		{
			name:              "success",
			requestUserID:     "user-id",
			claimsUserID:      "user-id",
			wantAccessToken:   "new-access-token",
			wantGetCalls:      1,
			wantGenerateCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.PatchRun(func() {
				var validatedToken, validatedType string
				validateMock := mockey.Mock(jwtutil.ValidateToken).
					To(func(token, tokenType string) (*jwtutil.Claims, error) {
						validatedToken = token
						validatedType = tokenType
						if tt.validateErr != nil {
							return nil, tt.validateErr
						}
						return &jwtutil.Claims{RegisteredClaims: gojwt.RegisteredClaims{
							Subject: tt.claimsUserID,
						}}, nil
					}).
					Build()

				var queriedToken string
				getMock := mockey.Mock((*repository.AuthDB).GetRefreshToken).
					To(func(_ *repository.AuthDB, _ context.Context, refreshToken string) (string, error) {
						queriedToken = refreshToken
						if tt.getErr != nil {
							return "", tt.getErr
						}
						return "1", nil
					}).
					Build()

				var generatedUserID, generatedType string
				generateMock := mockey.Mock(jwtutil.GenerateToken).
					To(func(userID, tokenType string) (string, error) {
						generatedUserID = userID
						generatedType = tokenType
						return "new-access-token", tt.generateErr
					}).
					Build()

				service := NewAuthService(&repository.AuthDB{})
				accessToken, err := service.RefreshToken(context.Background(), tt.requestUserID, "refresh-token")

				assertErrorCode(t, err, tt.wantCode)
				if tt.wantCause != nil && !errors.Is(err, tt.wantCause) {
					t.Fatalf("RefreshToken() error = %v, want cause %v", err, tt.wantCause)
				}
				if accessToken != tt.wantAccessToken {
					t.Fatalf("RefreshToken() access token = %q, want %q", accessToken, tt.wantAccessToken)
				}
				if validateMock.Times() != 1 {
					t.Fatalf("ValidateToken() calls = %d, want 1", validateMock.Times())
				}
				if validatedToken != "refresh-token" || validatedType != constants.TypeRefreshToken {
					t.Fatalf(
						"ValidateToken() arguments = (%q, %q), want (refresh-token, %q)",
						validatedToken,
						validatedType,
						constants.TypeRefreshToken,
					)
				}
				if getMock.Times() != tt.wantGetCalls {
					t.Fatalf("GetRefreshToken() calls = %d, want %d", getMock.Times(), tt.wantGetCalls)
				}
				if tt.wantGetCalls == 1 && queriedToken != "refresh-token" {
					t.Fatalf("GetRefreshToken() token = %q, want refresh-token", queriedToken)
				}
				if generateMock.Times() != tt.wantGenerateCalls {
					t.Fatalf("GenerateToken() calls = %d, want %d", generateMock.Times(), tt.wantGenerateCalls)
				}
				if tt.wantGenerateCalls == 1 {
					if generatedUserID != tt.requestUserID || generatedType != constants.TypeAccessToken {
						t.Fatalf(
							"GenerateToken() arguments = (%q, %q), want (%q, %q)",
							generatedUserID,
							generatedType,
							tt.requestUserID,
							constants.TypeAccessToken,
						)
					}
				}
			})
		})
	}
}
