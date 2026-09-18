package service

import (
	"context"
	"errors"
	"unicode"

	"github.com/Sna-ken/videoweb/internal/auth/repository"
	"github.com/Sna-ken/videoweb/pkg/constans"
	"github.com/Sna-ken/videoweb/pkg/db/model"
	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errno.Wrap(errno.ParamEmptyError, nil)
	}

	if valid, err := checkPasswordValid(password); !valid {
		return "", errno.Wrap(errno.NewErr(errno.PasswordInvalidErrorCode, "密码无效"), err)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", errno.Wrap(errno.NewErr(errno.PasswordHashErrorCode, "加密失败"), err)
	}

	exists, err := s.authdb.ExistsByUsername(ctx, username)
	if err != nil {
		return "", errno.Wrap(errno.InternalDatabaseError, err)
	}
	if exists {
		return "", errno.Wrap(
			errno.NewErr(errno.UserHasExistedErrorCode, "用户已存在"), nil)
	}

	userID := uuid.NewString()
	account := &model.AuthAccount{
		UserID:   userID,
		Username: username,
		Password: hashedPassword,
	}

	err = s.authdb.CreateAccount(ctx, account)
	switch {
	case err == nil:
		return userID, nil
	case errors.Is(err, repository.ErrAccountAlreadyExists):
		return "", errno.Wrap(
			errno.NewErr(errno.UserHasExistedErrorCode, "用户已存在"),
			err,
		)
	default:
		return "", errno.Wrap(errno.InternalDatabaseError, err)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func checkPasswordValid(password string) (bool, error) {
	if len(password) < constans.PasswordMinLength {
		return false, errors.New("密码长度小于6")
	}
	if len(password) > constans.PasswordMaxLength {
		return false, errors.New("密码长度大于20")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsSpace(r):
			return false, errors.New("密码中含有空格")
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	categoryCount := 0
	for _, matched := range []bool{hasUpper, hasLower, hasDigit, hasSpecial} {
		if matched {
			categoryCount++
		}
	}
	if categoryCount < constans.PasswordLeastCategoryCount {
		return false, errors.New("密码中必须至少含有大小写字母、数字、特殊符号中任意三种")
	}

	return true, nil
}
