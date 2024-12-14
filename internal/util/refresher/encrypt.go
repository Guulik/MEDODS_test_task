package refresher

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(refreshTokenRaw string) ([]byte, error) {
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshTokenRaw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return refreshHash, nil
}
