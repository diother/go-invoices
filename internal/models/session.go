package models

type Session struct {
	SessionToken int64 `db:"session_token"`
	UserID       int64 `db:"user_id"`
	ExpiresAt    int64 `db:"expires_at"`
}

func NewSession(sessionToken, userID, expiresAt int64) *Session {
	return &Session{
		SessionToken: sessionToken,
		UserID:       userID,
		ExpiresAt:    expiresAt,
	}
}
