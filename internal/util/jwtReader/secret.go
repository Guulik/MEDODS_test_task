package jwtReader

import (
	"bytes"
	"os"
)

func LoadJWTSecret() []byte {
	secretPath := "/run/secrets/jwt_secret"
	secret, err := os.ReadFile(secretPath)
	if err != nil {
		//TODO: delete me
		secret = bytes.NewBufferString("spy").Bytes()
		//TODO: uncomment me
		//log.Fatalf("Failed to load JWT secret: %v", err)
	}
	return secret
}
