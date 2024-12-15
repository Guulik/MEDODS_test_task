package service

import (
	"MEDODS-test/internal/configure"
	"MEDODS-test/internal/domain/model"
	"context"
	"log/slog"
	"time"
)

type Service struct {
	cfg           *configure.Config
	log           *slog.Logger
	tokenModifier TokenModifier
	tokenProvider TokenProvider
}

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=TokenProvider
type TokenProvider interface {
	Get(ctx context.Context, userID string) (*model.RefreshTokenDB, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=TokenModifier
type TokenModifier interface {
	Insert(ctx context.Context, userID string, tokenHash string, ipAddress string, expiresAt time.Time) error
	Delete(ctx context.Context, userID string) error
}

func New(
	cfg *configure.Config,
	log *slog.Logger,
	tokenModifier TokenModifier,
	tokenProvider TokenProvider,
) *Service {
	return &Service{
		cfg:           cfg,
		log:           log,
		tokenModifier: tokenModifier,
		tokenProvider: tokenProvider,
	}
}
