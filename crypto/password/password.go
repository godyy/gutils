package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Encrypt 加密密码
func Encrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verify 验证密码
// hashed 是通过 Encrypt 加密后的密码
func Verify(encrypted, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
}
