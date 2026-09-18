package user

import (
	"context"

	userservice "github.com/Sna-ken/videoweb/internal/user/service"
	user "github.com/Sna-ken/videoweb/kitex_gen/user"
	"github.com/Sna-ken/videoweb/pkg/base"
)

type userService interface {
	CreateUser(ctx context.Context, userID, username string) error
}

// UserServiceImpl implements the service defined in the IDL.
type UserServiceImpl struct {
	service userService
}

func NewUserServiceImpl(service *userservice.UserService) *UserServiceImpl {
	return &UserServiceImpl{service: service}
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {
	// TODO: Your code here...
	return
}

// UploadAvatar implements the UserServiceImpl interface.
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, req *user.UploadAvatarReq) (resp *user.UploadAvatarResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateSignature implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateSignature(ctx context.Context, req *user.UpdateSignatureReq) (resp *user.UpdateSignatureResp, err error) {
	// TODO: Your code here...
	return
}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserReq) (resp *user.CreateUserResp, err error) {
	var userID, username string
	if req.UserId != "" && req.Username != "" {
		userID = req.UserId
		username = req.Username
	}

	bizErr := s.service.CreateUser(ctx, userID, username)
	return &user.CreateUserResp{Base: base.BaseRPCResp(bizErr)}, nil
}
