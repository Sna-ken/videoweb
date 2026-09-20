package utils

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
