package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sna-ken/videoweb/internal/auth/repository"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/bytedance/mockey"
)

func TestLogout(t *testing.T) {
	databaseErr := errors.New("redis unavailable")

	tests := []struct {
		name            string
		refreshToken    string
		deleteErr       error
		wantCode        int64
		wantCause       error
		wantDeleteCalls int
	}{
		{
			name:     "empty refresh token",
			wantCode: errno.ParamEmptyErrorCode,
		},
		{
			name:            "delete refresh token failed",
			refreshToken:    "refresh-token",
			deleteErr:       databaseErr,
			wantCode:        errno.InternalDatabaseErrorCode,
			wantCause:       databaseErr,
			wantDeleteCalls: 1,
		},
		{
			name:            "success",
			refreshToken:    "refresh-token",
			wantDeleteCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.PatchRun(func() {
				var deletedToken string
				deleteMock := mockey.Mock((*repository.AuthDB).DeleteRefreshToken).
					To(func(_ *repository.AuthDB, _ context.Context, refreshToken string) error {
						deletedToken = refreshToken
						return tt.deleteErr
					}).
					Build()

				service := NewAuthService(&repository.AuthDB{})
				err := service.Logout(context.Background(), tt.refreshToken)

				assertErrorCode(t, err, tt.wantCode)
				if tt.wantCause != nil && !errors.Is(err, tt.wantCause) {
					t.Fatalf("Logout() error = %v, want cause %v", err, tt.wantCause)
				}
				if deleteMock.Times() != tt.wantDeleteCalls {
					t.Fatalf(
						"DeleteRefreshToken() calls = %d, want %d",
						deleteMock.Times(),
						tt.wantDeleteCalls,
					)
				}
				if tt.wantDeleteCalls == 1 && deletedToken != tt.refreshToken {
					t.Fatalf(
						"DeleteRefreshToken() argument = %q, want %q",
						deletedToken,
						tt.refreshToken,
					)
				}
			})
		})
	}
}
