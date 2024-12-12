package model

import (
	"time"
)

type RefreshTokenDB struct {
	UserID    string
	TokenHash string
	IPAddress string
	CreatedAt time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
