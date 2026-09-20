package auth

import (
	"context"
	"errors"

	authservice "github.com/Sna-ken/videoweb/internal/auth/service"
	auth "github.com/Sna-ken/videoweb/kitex_gen/auth"
	user "github.com/Sna-ken/videoweb/kitex_gen/user"
	"github.com/Sna-ken/videoweb/pkg/base"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/cloudwego/kitex/client/callopt"
)

type registerService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password, mfa_code string) (string, string, error)
}

type userCreator interface {
	CreateUser(ctx context.Context, req *user.CreateUserReq, callOptions ...callopt.Option) (*user.CreateUserResp, error)
}

// AuthServiceImpl implements the service defined in the IDL.
type AuthServiceImpl struct {
	service    registerService
	userClient userCreator
}

func NewAuthServiceImpl(service *authservice.AuthService, userClient userCreator) *AuthServiceImpl {
	return &AuthServiceImpl{
		service:    service,
		userClient: userClient,
	}
}

// Register implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Register(ctx context.Context, req *auth.RegisterReq) (resp *auth.RegisterResp, err error) {
	var username, password string
	if req.Password != "" && req.Username != "" {
		username = req.Username
		password = req.Password
	}

	userID, err := s.service.Register(ctx, username, password)
	if err != nil {
		return &auth.RegisterResp{Base: base.ErrsRPCResp(err)}, err
	}
	err = s.createUser(ctx, userID, username)
	if err != nil {
		return &auth.RegisterResp{Base: base.ErrsRPCResp(err)}, err
	}

	return &auth.RegisterResp{Base: base.SuccessRPCResp()}, nil
}

func (s *AuthServiceImpl) createUser(ctx context.Context, userID, username string) error {
	resp, err := s.userClient.CreateUser(ctx, &user.CreateUserReq{
		UserId:   userID,
		Username: username,
	})
	if err != nil {
		return errno.Wrap(errno.InternalServiceError, err)
	}
	if resp == nil || resp.Base == nil {
		return errno.Wrap(errno.InternalServiceError, errors.New("user service returned an empty response"))
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.NewErr(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

// Login implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Login(ctx context.Context, req *auth.LoginReq) (resp *auth.LoginResp, err error) {
	var username, password string
	if req.Password != "" && req.Username != "" {
		username = req.Username
		password = req.Password
	}

	accessToken, refreshToken, err := s.service.Login(ctx, username, password, req.MfaCode)
	if err != nil {
		return &auth.LoginResp{Base: base.ErrsRPCResp(err)}, err
	}

	return &auth.LoginResp{
		Base:         base.SuccessRPCResp(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Logout(ctx context.Context, req *auth.LogoutReq) (resp *auth.LogoutResp, err error) {
	// TODO: Your code here...
	return
}

// GetMFAqr implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) GetMFAqr(ctx context.Context, req *auth.GetMFAqrReq) (resp *auth.GetMFAqrResp, err error) {
	// TODO: Your code here...
	return
}

// BindMFA implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) BindMFA(ctx context.Context, req *auth.BindMFAReq) (resp *auth.BindMFAResp, err error) {
	// TODO: Your code here...
	return
}

// UnbindMFA implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) UnbindMFA(ctx context.Context, req *auth.UnbindMFAReq) (resp *auth.UnbindMFAResp, err error) {
	// TODO: Your code here...
	return
}

// RefreshToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *auth.RefreshTokenReq) (resp *auth.RefreshTokenResp, err error) {
	// TODO: Your code here...
	return
}
