package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sna-ken/videoweb/internal/auth/repository"
	"github.com/Sna-ken/videoweb/pkg/db/model"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/bytedance/mockey"
	"github.com/google/uuid"
)

func TestRegisterRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantCode int64
	}{
		{name: "empty username", password: "Abc123!", wantCode: errno.ParamEmptyErrorCode},
		{name: "empty password", username: "alice", wantCode: errno.ParamEmptyErrorCode},
		{name: "invalid password", username: "alice", password: "abc", wantCode: errno.PasswordInvalidErrorCode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAuthService(nil)
			_, err := service.Register(context.Background(), tt.username, tt.password)
			assertErrorCode(t, err, tt.wantCode)
		})
	}
}

func TestRegisterReturnsPasswordHashError(t *testing.T) {
	mockey.PatchRun(func() {
		hashErr := errors.New("hash failed")
		mockey.Mock((*repository.AuthDB).ExistsByUsername).Return(false, nil).Build()
		mockey.Mock(hashPassword).Return("", hashErr).Build()

		service := NewAuthService(&repository.AuthDB{})
		_, err := service.Register(context.Background(), "alice", "Abc123!")

		assertErrorCode(t, err, errno.PasswordHashErrorCode)
		if !errors.Is(err, hashErr) {
			t.Fatalf("Register() error = %v, want cause %v", err, hashErr)
		}
	})
}

func TestRegisterRepositoryResults(t *testing.T) {
	databaseErr := errors.New("database unavailable")
	tests := []struct {
		name            string
		exists          bool
		existsErr       error
		createErr       error
		wantCode        int64
		wantCause       error
		wantCreateCalls int
	}{
		{name: "success", wantCreateCalls: 1},
		{
			name:            "username query failed",
			existsErr:       databaseErr,
			wantCode:        errno.InternalDatabaseErrorCode,
			wantCause:       databaseErr,
			wantCreateCalls: 0,
		},
		{
			name:            "username already exists",
			exists:          true,
			wantCode:        errno.UserHasExistedErrorCode,
			wantCreateCalls: 0,
		},
		{
			name:            "duplicate detected while creating account",
			createErr:       repository.ErrAccountAlreadyExists,
			wantCode:        errno.UserHasExistedErrorCode,
			wantCause:       repository.ErrAccountAlreadyExists,
			wantCreateCalls: 1,
		},
		{
			name:            "create account failed",
			createErr:       databaseErr,
			wantCode:        errno.InternalDatabaseErrorCode,
			wantCause:       databaseErr,
			wantCreateCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.PatchRun(func() {
				mockey.Mock(hashPassword).Return("hashed-password", nil).Build()
				existsMock := mockey.Mock((*repository.AuthDB).ExistsByUsername).
					Return(tt.exists, tt.existsErr).
					Build()

				var createdAccount *model.AuthAccount
				createMock := mockey.Mock((*repository.AuthDB).CreateAccount).
					To(func(_ *repository.AuthDB, _ context.Context, account *model.AuthAccount) error {
						createdAccount = account
						return tt.createErr
					}).
					Build()

				service := NewAuthService(&repository.AuthDB{})
				userID, err := service.Register(context.Background(), "alice", "Abc123!")

				assertErrorCode(t, err, tt.wantCode)
				if tt.wantCause != nil && !errors.Is(err, tt.wantCause) {
					t.Fatalf("Register() error = %v, want cause %v", err, tt.wantCause)
				}
				if existsMock.Times() != 1 {
					t.Fatalf("ExistsByUsername() calls = %d, want 1", existsMock.Times())
				}
				if createMock.Times() != tt.wantCreateCalls {
					t.Fatalf("CreateAccount() calls = %d, want %d", createMock.Times(), tt.wantCreateCalls)
				}
				if tt.wantCode == 0 {
					assertCreatedAccount(t, createdAccount)
					if userID != createdAccount.UserID {
						t.Fatalf("returned user ID = %q, want %q", userID, createdAccount.UserID)
					}
				}
			})
		})
	}
}

func assertErrorCode(t *testing.T, err error, wantCode int64) {
	t.Helper()
	if wantCode == 0 {
		if err != nil {
			t.Fatalf("error = %v, want nil", err)
		}
		return
	}
	if err == nil {
		t.Fatalf("error = nil, want code %d", wantCode)
	}
	if got := errno.Convert(err).Code(); got != wantCode {
		t.Fatalf("error code = %d, want %d", got, wantCode)
	}
}

func assertCreatedAccount(t *testing.T, account *model.AuthAccount) {
	t.Helper()
	if account == nil {
		t.Fatal("CreateAccount() received nil account")
	}
	if account.Username != "alice" {
		t.Fatalf("created username = %q, want alice", account.Username)
	}
	if account.Password != "hashed-password" {
		t.Fatalf("created password = %q, want mocked hash", account.Password)
	}
	if _, err := uuid.Parse(account.UserID); err != nil {
		t.Fatalf("created user ID %q is not a valid UUID: %v", account.UserID, err)
	}
}
