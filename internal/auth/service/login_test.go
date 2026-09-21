package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sna-ken/videoweb/internal/auth/repository"
	"github.com/Sna-ken/videoweb/pkg/constants"
	"github.com/Sna-ken/videoweb/pkg/db/model"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/Sna-ken/videoweb/pkg/utils"
	"github.com/bytedance/mockey"
)

func TestLoginRejectsInvalidCredentialsAndMFA(t *testing.T) {
	databaseErr := errors.New("database unavailable")

	tests := []struct {
		name              string
		account           *model.AuthAccount
		searchErr         error
		passwordMatches   bool
		mfaCode           string
		mfaValid          bool
		wantCode          int64
		wantCause         error
		wantPasswordCalls int
		wantMFACalls      int
	}{
		{
			name:      "user not found",
			searchErr: databaseErr,
			wantCode:  errno.UserNotFoundErrorCode,
			wantCause: databaseErr,
		},
		{
			name:              "incorrect password",
			account:           loginAccount(false),
			wantCode:          errno.PasswordIncorrectErrorCode,
			wantPasswordCalls: 1,
		},
		{
			name:              "missing MFA code",
			account:           loginAccount(true),
			passwordMatches:   true,
			wantCode:          errno.MFAcodeEmptyErrorCode,
			wantPasswordCalls: 1,
		},
		{
			name:              "incorrect MFA code",
			account:           loginAccount(true),
			passwordMatches:   true,
			mfaCode:           "123456",
			wantCode:          errno.MFAcodeErrorCode,
			wantPasswordCalls: 1,
			wantMFACalls:      1,
		},
		{
			name:              "MFA code supplied when MFA is disabled",
			account:           loginAccount(false),
			passwordMatches:   true,
			mfaCode:           "123456",
			wantCode:          errno.MFANotEnabledErrorCode,
			wantPasswordCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.PatchRun(func() {
				searchMock := mockey.Mock((*repository.AuthDB).SearchByName).
					Return(tt.account, tt.searchErr).
					Build()
				passwordMock := mockey.Mock(utils.CheckPasswordHash).
					Return(tt.passwordMatches).
					Build()
				mfaMock := mockey.Mock(utils.ValidateMFA).
					Return(tt.mfaValid).
					Build()
				tokenMock := mockey.Mock(jwt.GenerateToken).
					Return("unused-token", nil).
					Build()
				saveMock := mockey.Mock((*repository.AuthDB).SaveRefreshToken).
					Return(nil).
					Build()

				service := NewAuthService(&repository.AuthDB{})
				accessToken, refreshToken, err := service.Login(
					context.Background(),
					"alice",
					"Abc123!",
					tt.mfaCode,
				)

				assertErrorCode(t, err, tt.wantCode)
				if tt.wantCause != nil && !errors.Is(err, tt.wantCause) {
					t.Fatalf("Login() error = %v, want cause %v", err, tt.wantCause)
				}
				if accessToken != "" || refreshToken != "" {
					t.Fatalf("Login() tokens = (%q, %q), want empty tokens", accessToken, refreshToken)
				}
				if searchMock.Times() != 1 {
					t.Fatalf("SearchByName() calls = %d, want 1", searchMock.Times())
				}
				if passwordMock.Times() != tt.wantPasswordCalls {
					t.Fatalf("CheckPasswordHash() calls = %d, want %d", passwordMock.Times(), tt.wantPasswordCalls)
				}
				if mfaMock.Times() != tt.wantMFACalls {
					t.Fatalf("ValidateMFA() calls = %d, want %d", mfaMock.Times(), tt.wantMFACalls)
				}
				if tokenMock.Times() != 0 {
					t.Fatalf("GenerateToken() calls = %d, want 0", tokenMock.Times())
				}
				if saveMock.Times() != 0 {
					t.Fatalf("SaveRefreshToken() calls = %d, want 0", saveMock.Times())
				}
			})
		})
	}
}

func TestLoginGeneratesAndPersistsTokens(t *testing.T) {
	accessTokenErr := errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "access token generation failed"), nil)
	refreshTokenErr := errno.Wrap(errno.NewErr(errno.InternalServiceErrorCode, "refresh token generation failed"), nil)
	databaseErr := errors.New("database unavailable")

	tests := []struct {
		name            string
		mfaEnabled      bool
		mfaCode         string
		accessTokenErr  error
		refreshTokenErr error
		saveErr         error
		wantCode        int64
		wantCause       error
		wantMFACalls    int
		wantTokenCalls  int
		wantSaveCalls   int
	}{
		{
			name:           "success without MFA",
			wantTokenCalls: 2,
			wantSaveCalls:  1,
		},
		{
			name:           "success with MFA",
			mfaEnabled:     true,
			mfaCode:        "123456",
			wantMFACalls:   1,
			wantTokenCalls: 2,
			wantSaveCalls:  1,
		},
		{
			name:           "access token generation failed",
			accessTokenErr: accessTokenErr,
			wantCode:       errno.InternalServiceErrorCode,
			wantCause:      accessTokenErr,
			wantTokenCalls: 1,
		},
		{
			name:            "refresh token generation failed",
			refreshTokenErr: refreshTokenErr,
			wantCode:        errno.InternalServiceErrorCode,
			wantCause:       refreshTokenErr,
			wantTokenCalls:  2,
		},
		{
			name:           "refresh token persistence failed",
			saveErr:        databaseErr,
			wantCode:       errno.InternalDatabaseErrorCode,
			wantCause:      databaseErr,
			wantTokenCalls: 2,
			wantSaveCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.PatchRun(func() {
				account := loginAccount(tt.mfaEnabled)
				searchMock := mockey.Mock((*repository.AuthDB).SearchByName).
					Return(account, nil).
					Build()
				passwordMock := mockey.Mock(utils.CheckPasswordHash).
					Return(true).
					Build()
				mfaMock := mockey.Mock(utils.ValidateMFA).
					Return(true).
					Build()
				var generatedUserIDs, generatedTypes []string
				tokenMock := mockey.Mock(jwt.GenerateToken).
					To(func(userID, tokenType string) (string, error) {
						generatedUserIDs = append(generatedUserIDs, userID)
						generatedTypes = append(generatedTypes, tokenType)
						switch tokenType {
						case constants.TypeAccessToken:
							return "access-token", tt.accessTokenErr
						case constants.TypeRefreshToken:
							return "refresh-token", tt.refreshTokenErr
						default:
							t.Fatalf("GenerateToken() token type = %q, want access_token or refresh_token", tokenType)
							return "", nil
						}
					}).
					Build()

				var savedToken string
				saveMock := mockey.Mock((*repository.AuthDB).SaveRefreshToken).
					To(func(_ *repository.AuthDB, _ context.Context, refreshToken string) error {
						savedToken = refreshToken
						return tt.saveErr
					}).
					Build()

				service := NewAuthService(&repository.AuthDB{})
				accessToken, refreshToken, err := service.Login(
					context.Background(),
					"alice",
					"Abc123!",
					tt.mfaCode,
				)

				assertErrorCode(t, err, tt.wantCode)
				if tt.wantCause != nil && !errors.Is(err, tt.wantCause) {
					t.Fatalf("Login() error = %v, want cause %v", err, tt.wantCause)
				}
				if searchMock.Times() != 1 || passwordMock.Times() != 1 {
					t.Fatalf(
						"dependency calls = search:%d password:%d, want 1 each",
						searchMock.Times(),
						passwordMock.Times(),
					)
				}
				if tokenMock.Times() != tt.wantTokenCalls {
					t.Fatalf("GenerateToken() calls = %d, want %d", tokenMock.Times(), tt.wantTokenCalls)
				}
				wantTypes := []string{constants.TypeAccessToken, constants.TypeRefreshToken}
				for i := 0; i < tt.wantTokenCalls; i++ {
					if generatedUserIDs[i] != "user-id" || generatedTypes[i] != wantTypes[i] {
						t.Fatalf(
							"GenerateToken() call %d = (%q, %q), want (user-id, %q)",
							i+1,
							generatedUserIDs[i],
							generatedTypes[i],
							wantTypes[i],
						)
					}
				}
				if mfaMock.Times() != tt.wantMFACalls {
					t.Fatalf("ValidateMFA() calls = %d, want %d", mfaMock.Times(), tt.wantMFACalls)
				}
				if saveMock.Times() != tt.wantSaveCalls {
					t.Fatalf("SaveRefreshToken() calls = %d, want %d", saveMock.Times(), tt.wantSaveCalls)
				}

				if tt.wantCode == 0 {
					if accessToken != "access-token" || refreshToken != "refresh-token" {
						t.Fatalf("Login() tokens = (%q, %q), want access-token and refresh-token", accessToken, refreshToken)
					}
					if savedToken != "refresh-token" {
						t.Fatalf("saved refresh token = %q, want refresh-token", savedToken)
					}
					return
				}
				if accessToken != "" || refreshToken != "" {
					t.Fatalf("Login() tokens = (%q, %q), want empty tokens", accessToken, refreshToken)
				}
			})
		})
	}
}

func loginAccount(mfaEnabled bool) *model.AuthAccount {
	return &model.AuthAccount{
		UserID:     "user-id",
		Username:   "alice",
		Password:   "hashed-password",
		MFAEnabled: mfaEnabled,
		MFASecret:  "mfa-secret",
	}
}
