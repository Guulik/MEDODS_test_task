create table if not exists "user"(
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL
);


create table if not exists "refresh_token"
(
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
