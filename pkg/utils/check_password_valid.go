package utils

import (
	"errors"
	"unicode"

	"github.com/Sna-ken/videoweb/pkg/constants"
)

func CheckPasswordValid(password string) (bool, error) {
	if len(password) < constants.PasswordMinLength {
		return false, errors.New("密码长度小于6")
	}
	if len(password) > constants.PasswordMaxLength {
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
	if categoryCount < constants.PasswordLeastCategoryCount {
		return false, errors.New("密码中必须至少含有大小写字母、数字、特殊符号中任意三种")
	}

	return true, nil
}
