package service

import (
	"MEDODS-test/internal/configure"
	"MEDODS-test/internal/domain/model"
	"context"
	"log/slog"
	"time"
)

type Service struct {
	jwtSecret []byte

	cfg           *configure.Config
	log           *slog.Logger
	tokenModifier TokenModifier
	tokenProvider TokenProvider
}

type TokenProvider interface {
	Get(ctx context.Context, userID string) (*model.RefreshTokenDB, error)
}
type TokenModifier interface {
	Insert(ctx context.Context, userID string, tokenHash string, ipAddress string, expiresAt time.Time) error
	Delete(ctx context.Context, userID string) error
}

func New(
	jwtSecret []byte,
	cfg *configure.Config,
	log *slog.Logger,
	tokenModifier TokenModifier,
	tokenProvider TokenProvider,
) *Service {
	return &Service{
		jwtSecret:     jwtSecret,
		cfg:           cfg,
		log:           log,
		tokenModifier: tokenModifier,
		tokenProvider: tokenProvider,
	}
}
