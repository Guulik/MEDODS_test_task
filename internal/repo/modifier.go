package repo

import "context"

//TODO: изменение даты использования токена.

// TODO: проверить
func (r *Storage) Insert(ctx context.Context, userID string, tokenHash string, ipAddress string) error {
	query := `INSERT INTO refresh_token (user_id, refresh_token_hash, ip_address) 
              VALUES ($1, $2, $3)`

	_, err := r.pool.Exec(ctx, query, userID, tokenHash, ipAddress)
	return err
}

// TODO: проверить
func (r *Storage) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM refresh_token WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
