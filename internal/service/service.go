package service

import (
	"MEDODS-test/internal/domain/model"
	"context"
)

type Service struct {
	Tokens Tokens
	Email  EmailNotifier
}

type Tokens interface {
	GenerateTokens(ctx context.Context, userID, ipAddress string) (*model.TokenPair, error)
	RefreshTokens(ctx context.Context, userID, refreshTokenRaw, ipAddress string) (*model.TokenPair, error)
}

type EmailNotifier interface {
	SendWarning(ctx context.Context, address string, ip string) error
}

func New(
	Tokens Tokens,
	Email EmailNotifier,
) *Service {
	return &Service{
		Tokens: Tokens,
		Email:  Email,
	}
}
