package repo

import (
	"MEDODS-test/internal/domain/model"
	"context"
)

func (r *Storage) Get(ctx context.Context, userID string) (*model.RefreshTokenDB, error) {
	query := `SELECT refresh_token_hash, ip_address, created_at 
              FROM refresh_tokens WHERE user_id = $1`
	row := r.pool.QueryRow(ctx, query, userID)

	var token model.RefreshTokenDB
	token.UserID = userID
	if err := row.Scan(&token.TokenHash, &token.IPAddress, &token.CreatedAt); err != nil {
		return nil, err
	}

	return &token, nil
}
