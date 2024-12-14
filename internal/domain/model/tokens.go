package model

import (
	"time"
)

type RefreshTokenDB struct {
	UserGuid  string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	IPAddress string    `db:"ip_address"`
	ExpiresAt time.Time `db:"expires_at"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
