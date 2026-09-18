package service

import "github.com/Sna-ken/videoweb/internal/auth/repository"

type AuthService struct {
	authdb *repository.AuthDB
}

func NewAuthService(authdb *repository.AuthDB) *AuthService {
	return &AuthService{authdb: authdb}
}
