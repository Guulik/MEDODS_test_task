package jwtReader

import (
	"log"
	"os"
)

func LoadJWTSecret() []byte {
	secretPath := "/run/secrets/jwt"
	secret, err := os.ReadFile(secretPath)
	if err != nil {
		log.Fatalf("Failed to load JWT secret: %v", err)
	}
	return secret
}
