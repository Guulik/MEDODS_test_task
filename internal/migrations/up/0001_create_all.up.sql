create table if not exists "refresh_token"
(
    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
