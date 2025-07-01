package domain

import "time"

type UserSession struct {
	Token     string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewUserSession(userID string, expiresAt time.Time) *UserSession {
	return &UserSession{
		Token:     "",
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
}
