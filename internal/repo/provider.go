package repo

import (
	"MEDODS-test/internal/domain/model"
	"context"
)

func (r *Storage) Get(ctx context.Context, userID string) (*model.RefreshTokenDB, error) {
	var token model.RefreshTokenDB
	token.UserGuid = userID

	query := `SELECT token_hash, ip_address, expires_at 
              FROM refresh_token WHERE user_id = $1`

	row := r.pool.QueryRow(ctx, query, userID)

	if err := row.Scan(&token.TokenHash, &token.IPAddress, &token.ExpiresAt); err != nil {
		return nil, err
	}

	return &token, nil
}
