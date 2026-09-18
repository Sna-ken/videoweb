package service

import "github.com/Sna-ken/videoweb/internal/user/repository"

type UserService struct {
	userdb *repository.UserDB
}

func NewUserService(userdb *repository.UserDB) *UserService {
	return &UserService{userdb: userdb}
}
