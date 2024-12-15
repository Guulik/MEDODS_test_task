package jwtReader

import (
	"bytes"
	"log"
	"os"
	"strings"
)

func LoadJWTSecret(env string) []byte {
	if strings.EqualFold(env, "local") {
		return bytes.NewBufferString("default").Bytes()
	}

	secretPath := "/run/secrets/jwt"
	secret, err := os.ReadFile(secretPath)
	if err != nil {
		log.Fatalf("Failed to load JWT secret: %v", err)
	}

	return secret
}
