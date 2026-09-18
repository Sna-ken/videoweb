package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sna-ken/videoweb/internal/user/repository"
	"github.com/Sna-ken/videoweb/pkg/db/model"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/bytedance/mockey"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCreateUserRejectsEmptyParams(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		username string
	}{
		{name: "empty user ID", username: "alice"},
		{name: "empty username", userID: "user-id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewUserService(nil)
			err := service.CreateUser(context.Background(), tt.userID, tt.username)
			assertUserErrorCode(t, err, errno.ParamEmptyErrorCode)
		})
	}
}

func TestCreateUserRepositoryResults(t *testing.T) {
	databaseErr := errors.New("database unavailable")
	tests := []struct {
		name      string
		createErr error
		wantCode  int64
		wantCause error
	}{
		{name: "success"},
		{
			name:      "duplicate user",
			createErr: gorm.ErrDuplicatedKey,
			wantCode:  errno.UserHasExistedErrorCode,
			wantCause: gorm.ErrDuplicatedKey,
		},
		{
			name:      "database failure",
			createErr: databaseErr,
			wantCode:  errno.InternalDatabaseErrorCode,
			wantCause: databaseErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockey.PatchRun(func() {
				var createdUser *model.User
				mockey.Mock((*repository.UserDB).CreateUser).
					To(func(_ *repository.UserDB, _ context.Context, user *model.User) error {
						createdUser = user
						return tt.createErr
					}).
					Build()

				service := NewUserService(&repository.UserDB{})
				err := service.CreateUser(context.Background(), "auth-user-id", "alice")

				assertUserErrorCode(t, err, tt.wantCode)
				if tt.wantCause != nil && !errors.Is(err, tt.wantCause) {
					t.Fatalf("CreateUser() error = %v, want cause %v", err, tt.wantCause)
				}
				assertCreatedUser(t, createdUser)
			})
		})
	}
}

func assertUserErrorCode(t *testing.T, err error, wantCode int64) {
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

func assertCreatedUser(t *testing.T, user *model.User) {
	t.Helper()
	if user == nil {
		t.Fatal("CreateUser() received nil user")
	}
	if user.UserID != "auth-user-id" {
		t.Fatalf("created user ID = %q, want auth-user-id", user.UserID)
	}
	if user.Nickname != "alice" {
		t.Fatalf("created nickname = %q, want alice", user.Nickname)
	}
	if _, err := uuid.Parse(user.ID); err != nil {
		t.Fatalf("created profile ID %q is not a valid UUID: %v", user.ID, err)
	}
}
