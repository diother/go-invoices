CREATE TABLE sessions (
    session_token BIGINT PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    expires_at BIGINT NOT NULL
);
