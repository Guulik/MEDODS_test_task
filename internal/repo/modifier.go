package repo

import (
	"context"
	"time"
)

func (r *Storage) Insert(ctx context.Context, userID string, tokenHash string, ipAddress string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_token (user_id, token_hash, ip_address, expires_at) 
              VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(ctx, query, userID, tokenHash, ipAddress, expiresAt)
	return err
}

func (r *Storage) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM refresh_token WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
