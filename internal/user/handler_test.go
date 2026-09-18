package user

import (
	"context"
	"testing"

	usermodel "github.com/Sna-ken/videoweb/kitex_gen/user"
	"github.com/Sna-ken/videoweb/pkg/errno"
)

type stubUserService struct {
	userID   string
	username string
	err      error
}

func (s *stubUserService) CreateUser(_ context.Context, userID, username string) error {
	s.userID = userID
	s.username = username
	return s.err
}

func TestCreateUserCallsServiceAndBuildsResponse(t *testing.T) {
	service := &stubUserService{}
	handler := &UserServiceImpl{service: service}

	resp, err := handler.CreateUser(context.Background(), &usermodel.CreateUserReq{
		UserId:   "user-id",
		Username: "alice",
	})

	if err != nil {
		t.Fatalf("CreateUser() transport error = %v, want nil", err)
	}
	if resp.Base.Code != errno.SuccessCode {
		t.Fatalf("CreateUser() code = %d, want %d", resp.Base.Code, errno.SuccessCode)
	}
	if service.userID != "user-id" || service.username != "alice" {
		t.Fatalf("service params = %q, %q, want user-id, alice", service.userID, service.username)
	}
}

func TestCreateUserBuildsBusinessErrorResponse(t *testing.T) {
	service := &stubUserService{err: errno.Wrap(errno.ParamEmptyError, nil)}
	handler := &UserServiceImpl{service: service}

	resp, err := handler.CreateUser(context.Background(), nil)

	if err != nil {
		t.Fatalf("CreateUser() transport error = %v, want nil", err)
	}
	if resp.Base.Code != errno.ParamEmptyErrorCode {
		t.Fatalf("CreateUser() code = %d, want %d", resp.Base.Code, errno.ParamEmptyErrorCode)
	}
}
