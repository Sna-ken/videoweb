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
	userID              string
	err                 error
	accessToken         string
	refreshToken        string
	loginErr            error
	loginCalls          int
	loginUsername       string
	loginPassword       string
	loginMFACode        string
	refreshAccessToken  string
	refreshErr          error
	refreshCalls        int
	refreshUserID       string
	refreshRequestToken string
	logoutErr           error
	logoutCalls         int
	logoutRefreshToken  string
}

func (s *stubRegisterService) Register(context.Context, string, string) (string, error) {
	return s.userID, s.err
}

func (s *stubRegisterService) Login(_ context.Context, username, password, mfaCode string) (string, string, error) {
	s.loginCalls++
	s.loginUsername = username
	s.loginPassword = password
	s.loginMFACode = mfaCode
	return s.accessToken, s.refreshToken, s.loginErr
}

func (s *stubRegisterService) RefreshToken(_ context.Context, userID, refreshToken string) (string, error) {
	s.refreshCalls++
	s.refreshUserID = userID
	s.refreshRequestToken = refreshToken
	return s.refreshAccessToken, s.refreshErr
}

func (s *stubRegisterService) Logout(_ context.Context, refreshToken string) error {
	s.logoutCalls++
	s.logoutRefreshToken = refreshToken
	return s.logoutErr
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
	registerErr := errno.Wrap(errno.ParamEmptyError, nil)
	handler := &AuthServiceImpl{
		service: &stubRegisterService{
			err: registerErr,
		},
		userClient: creator,
	}

	resp, err := handler.Register(context.Background(), &authmodel.RegisterReq{})

	if !errors.Is(err, registerErr) {
		t.Fatalf("Register() error = %v, want %v", err, registerErr)
	}
	if resp.Base.Code != errno.ParamEmptyErrorCode {
		t.Fatalf("Register() code = %d, want %d", resp.Base.Code, errno.ParamEmptyErrorCode)
	}
	if creator.calls != 0 {
		t.Fatalf("CreateUser() calls = %d, want 0", creator.calls)
	}
}

func TestRegisterWrapsUserRPCError(t *testing.T) {
	rpcErr := errors.New("rpc unavailable")
	creator := &stubUserCreator{err: rpcErr}
	handler := &AuthServiceImpl{
		service:    &stubRegisterService{userID: "user-id"},
		userClient: creator,
	}

	resp, err := handler.Register(context.Background(), &authmodel.RegisterReq{Username: "alice", Password: "Abc123!"})

	if !errors.Is(err, rpcErr) {
		t.Fatalf("Register() error = %v, want cause %v", err, rpcErr)
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

	if got := errno.Convert(err).Code(); got != errno.UserHasExistedErrorCode {
		t.Fatalf("Register() error code = %d, want %d", got, errno.UserHasExistedErrorCode)
	}
	if resp.Base.Code != errno.UserHasExistedErrorCode {
		t.Fatalf("Register() code = %d, want %d", resp.Base.Code, errno.UserHasExistedErrorCode)
	}
}

func TestLoginReturnsTokens(t *testing.T) {
	service := &stubRegisterService{
		accessToken:  "access-token",
		refreshToken: "refresh-token",
	}
	handler := &AuthServiceImpl{service: service}

	resp, err := handler.Login(context.Background(), &authmodel.LoginReq{
		Username: "alice",
		Password: "Abc123!",
		MfaCode:  "123456",
	})

	if err != nil {
		t.Fatalf("Login() error = %v, want nil", err)
	}
	if resp.Base.Code != errno.SuccessCode {
		t.Fatalf("Login() code = %d, want %d", resp.Base.Code, errno.SuccessCode)
	}
	if resp.AccessToken != "access-token" || resp.RefreshToken != "refresh-token" {
		t.Fatalf("Login() tokens = (%q, %q), want access-token and refresh-token", resp.AccessToken, resp.RefreshToken)
	}
	if service.loginCalls != 1 {
		t.Fatalf("service Login() calls = %d, want 1", service.loginCalls)
	}
	if service.loginUsername != "alice" || service.loginPassword != "Abc123!" || service.loginMFACode != "123456" {
		t.Fatalf(
			"service Login() arguments = (%q, %q, %q), want alice, Abc123!, 123456",
			service.loginUsername,
			service.loginPassword,
			service.loginMFACode,
		)
	}
}

func TestLoginReturnsServiceError(t *testing.T) {
	loginErr := errno.Wrap(errno.NewErr(errno.PasswordIncorrectErrorCode, "密码错误"), nil)
	service := &stubRegisterService{loginErr: loginErr}
	handler := &AuthServiceImpl{service: service}

	resp, err := handler.Login(context.Background(), &authmodel.LoginReq{
		Username: "alice",
		Password: "wrong-password",
	})

	if !errors.Is(err, loginErr) {
		t.Fatalf("Login() error = %v, want %v", err, loginErr)
	}
	if resp.Base.Code != errno.PasswordIncorrectErrorCode {
		t.Fatalf("Login() code = %d, want %d", resp.Base.Code, errno.PasswordIncorrectErrorCode)
	}
	if resp.AccessToken != "" || resp.RefreshToken != "" {
		t.Fatalf("Login() tokens = (%q, %q), want empty tokens", resp.AccessToken, resp.RefreshToken)
	}
}

func TestRefreshTokenReturnsAccessToken(t *testing.T) {
	service := &stubRegisterService{refreshAccessToken: "new-access-token"}
	handler := &AuthServiceImpl{service: service}

	resp, err := handler.RefreshToken(context.Background(), &authmodel.RefreshTokenReq{
		Id:           "user-id",
		RefreshToken: "refresh-token",
	})

	if err != nil {
		t.Fatalf("RefreshToken() error = %v, want nil", err)
	}
	if resp.Base.Code != errno.SuccessCode {
		t.Fatalf("RefreshToken() code = %d, want %d", resp.Base.Code, errno.SuccessCode)
	}
	if resp.AccessToken != "new-access-token" {
		t.Fatalf("RefreshToken() access token = %q, want new-access-token", resp.AccessToken)
	}
	if service.refreshCalls != 1 {
		t.Fatalf("service RefreshToken() calls = %d, want 1", service.refreshCalls)
	}
	if service.refreshUserID != "user-id" {
		t.Fatalf("service RefreshToken() user ID = %q, want user-id", service.refreshUserID)
	}
	if service.refreshRequestToken != "refresh-token" {
		t.Fatalf(
			"service RefreshToken() argument = %q, want refresh-token",
			service.refreshRequestToken,
		)
	}
}

func TestRefreshTokenReturnsServiceError(t *testing.T) {
	refreshErr := errno.Wrap(errno.AuthInvalidError, errors.New("refresh token invalid"))
	service := &stubRegisterService{refreshErr: refreshErr}
	handler := &AuthServiceImpl{service: service}

	resp, err := handler.RefreshToken(context.Background(), &authmodel.RefreshTokenReq{
		Id:           "user-id",
		RefreshToken: "invalid-refresh-token",
	})

	if !errors.Is(err, refreshErr) {
		t.Fatalf("RefreshToken() error = %v, want %v", err, refreshErr)
	}
	if resp.Base.Code != errno.AuthInvalidErrorCode {
		t.Fatalf("RefreshToken() code = %d, want %d", resp.Base.Code, errno.AuthInvalidErrorCode)
	}
	if resp.AccessToken != "" {
		t.Fatalf("RefreshToken() access token = %q, want empty", resp.AccessToken)
	}
}

func TestLogoutReturnsSuccess(t *testing.T) {
	service := &stubRegisterService{}
	handler := &AuthServiceImpl{service: service}

	resp, err := handler.Logout(context.Background(), &authmodel.LogoutReq{
		RefreshToken: "refresh-token",
	})

	if err != nil {
		t.Fatalf("Logout() error = %v, want nil", err)
	}
	if resp.Base.Code != errno.SuccessCode {
		t.Fatalf("Logout() code = %d, want %d", resp.Base.Code, errno.SuccessCode)
	}
	if service.logoutCalls != 1 {
		t.Fatalf("service Logout() calls = %d, want 1", service.logoutCalls)
	}
	if service.logoutRefreshToken != "refresh-token" {
		t.Fatalf(
			"service Logout() argument = %q, want refresh-token",
			service.logoutRefreshToken,
		)
	}
}

func TestLogoutReturnsServiceError(t *testing.T) {
	logoutErr := errno.Wrap(errno.InternalDatabaseError, errors.New("redis unavailable"))
	service := &stubRegisterService{logoutErr: logoutErr}
	handler := &AuthServiceImpl{service: service}

	resp, err := handler.Logout(context.Background(), &authmodel.LogoutReq{
		RefreshToken: "refresh-token",
	})

	if !errors.Is(err, logoutErr) {
		t.Fatalf("Logout() error = %v, want %v", err, logoutErr)
	}
	if resp.Base.Code != errno.InternalDatabaseErrorCode {
		t.Fatalf(
			"Logout() code = %d, want %d",
			resp.Base.Code,
			errno.InternalDatabaseErrorCode,
		)
	}
	if service.logoutCalls != 1 {
		t.Fatalf("service Logout() calls = %d, want 1", service.logoutCalls)
	}
	if service.logoutRefreshToken != "refresh-token" {
		t.Fatalf(
			"service Logout() argument = %q, want refresh-token",
			service.logoutRefreshToken,
		)
	}
}
