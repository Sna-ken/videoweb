package utils

import "github.com/pquerna/otp/totp"

func ValidateMFA(code string, secret string) bool {
	return totp.Validate(code, secret)
}
