package service

import (
	"context"

	"github.com/Sna-ken/videoweb/pkg/errno"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/Sna-ken/videoweb/pkg/utils"
)

func (s *AuthService) Login(ctx context.Context, username, password, mfa_code string) (string, string, error) {
	info, err := s.authdb.SearchByName(ctx, username)
	if err != nil {
		return "", "", errno.Wrap(errno.NewErr(errno.UserNotFoundErrorCode, "用户不存在"), err)
	}

	if !utils.CheckPasswordHash(password, info.Password) {
		return "", "", errno.Wrap(errno.NewErr(errno.PasswordIncorrectErrorCode, "密码错误"), nil)
	}

	if info.MFAEnabled {
		if mfa_code == "" {
			return "", "", errno.Wrap(errno.NewErr(errno.MFAcodeEmptyErrorCode, "MFA验证码为空"), nil)
		}
		if !utils.ValidateMFA(mfa_code, info.MFASecret) {
			return "", "", errno.Wrap(errno.NewErr(errno.MFAcodeErrorCode, "MFA验证码错误"), nil)
		}
	}

	if !info.MFAEnabled && mfa_code != "" {
		return "", "", errno.Wrap(errno.NewErr(errno.MFANotEnabledErrorCode, "MFA未启用"), nil)
	}

	accessToken, refreshToken, jwterr := jwt.GenerateToken(info.UserID)
	if jwterr != nil {
		return "", "", jwterr
	}

	err = s.authdb.SaveRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errno.Wrap(errno.NewErr(errno.InternalDatabaseErrorCode, "保存refresh_token失败"), err)
	}

	return accessToken, refreshToken, nil
}
