CREATE TABLE sessions (
    session_token BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expires_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
