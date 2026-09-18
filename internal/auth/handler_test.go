package auth

import (
	"context"
	"errors"
	"testing"

	authmodel "github.com/Sna-ken/videoweb/kitex_gen/auth"
	"github.com/Sna-ken/videoweb/kitex_gen/model"
	user "github.com/Sna-ken/videoweb/kitex_gen/user"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/cloudwego/kitex/client/callopt"
)

type stubRegisterService struct {
	userID string
	err    error
}

func (s *stubRegisterService) Register(context.Context, string, string) (string, error) {
	return s.userID, s.err
}

type stubUserCreator struct {
	req   *user.CreateUserReq
	resp  *user.CreateUserResp
	err   error
	calls int
}

func (s *stubUserCreator) CreateUser(_ context.Context, req *user.CreateUserReq, _ ...callopt.Option) (*user.CreateUserResp, error) {
	s.calls++
	s.req = req
	return s.resp, s.err
}

func TestRegisterCreatesAuthAccountAndUserProfile(t *testing.T) {
	creator := &stubUserCreator{resp: &user.CreateUserResp{
		Base: &model.BaseResp{Code: errno.SuccessCode, Msg: errno.Success.Message()},
	}}
	handler := &AuthServiceImpl{
		service:    &stubRegisterService{userID: "user-id"},
		userClient: creator,
	}

	resp, err := handler.Register(context.Background(), &authmodel.RegisterReq{
		Username: "alice",
		Password: "Abc123!",
	})

	if err != nil {
		t.Fatalf("Register() transport error = %v, want nil", err)
	}
	if resp.Base.Code != errno.SuccessCode {
		t.Fatalf("Register() code = %d, want %d", resp.Base.Code, errno.SuccessCode)
	}
	if creator.calls != 1 {
		t.Fatalf("CreateUser() calls = %d, want 1", creator.calls)
	}
	if creator.req.UserId != "user-id" || creator.req.Username != "alice" {
		t.Fatalf("CreateUser() request = %+v, want user-id and alice", creator.req)
	}
}

func TestRegisterSkipsUserRPCWhenAuthRegistrationFails(t *testing.T) {
	creator := &stubUserCreator{}
	handler := &AuthServiceImpl{
		service: &stubRegisterService{
			err: errno.Wrap(errno.ParamEmptyError, nil),
		},
		userClient: creator,
	}

	resp, err := handler.Register(context.Background(), &authmodel.RegisterReq{})

	if err != nil {
		t.Fatalf("Register() transport error = %v, want nil", err)
	}
	if resp.Base.Code != errno.ParamEmptyErrorCode {
		t.Fatalf("Register() code = %d, want %d", resp.Base.Code, errno.ParamEmptyErrorCode)
	}
	if creator.calls != 0 {
		t.Fatalf("CreateUser() calls = %d, want 0", creator.calls)
	}
}

func TestRegisterWrapsUserRPCError(t *testing.T) {
	creator := &stubUserCreator{err: errors.New("rpc unavailable")}
	handler := &AuthServiceImpl{
		service:    &stubRegisterService{userID: "user-id"},
		userClient: creator,
	}

	resp, err := handler.Register(context.Background(), &authmodel.RegisterReq{Username: "alice", Password: "Abc123!"})

	if err != nil {
		t.Fatalf("Register() transport error = %v, want nil", err)
	}
	if resp.Base.Code != errno.InternalServiceErrorCode {
		t.Fatalf("Register() code = %d, want %d", resp.Base.Code, errno.InternalServiceErrorCode)
	}
}

func TestRegisterPropagatesUserBusinessError(t *testing.T) {
	creator := &stubUserCreator{resp: &user.CreateUserResp{
		Base: &model.BaseResp{Code: errno.UserHasExistedErrorCode, Msg: "用户已存在"},
	}}
	handler := &AuthServiceImpl{
		service:    &stubRegisterService{userID: "user-id"},
		userClient: creator,
	}

	resp, err := handler.Register(context.Background(), &authmodel.RegisterReq{Username: "alice", Password: "Abc123!"})

	if err != nil {
		t.Fatalf("Register() transport error = %v, want nil", err)
	}
	if resp.Base.Code != errno.UserHasExistedErrorCode {
		t.Fatalf("Register() code = %d, want %d", resp.Base.Code, errno.UserHasExistedErrorCode)
	}
}
