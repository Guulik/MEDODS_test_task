package repo

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

type Storage struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func New(
	log *slog.Logger,
	pool *pgxpool.Pool,
) *Storage {
	return &Storage{
		log:  log,
		pool: pool,
	}
}
