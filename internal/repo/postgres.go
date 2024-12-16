package repo

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
)

type PostgresStorage struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func New(
	log *slog.Logger,
	pool *pgxpool.Pool,
) *PostgresStorage {
	return &PostgresStorage{
		log:  log,
		pool: pool,
	}
}
