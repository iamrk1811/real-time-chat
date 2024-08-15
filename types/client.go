package types

import "time"

type User struct {
	Username string
	Password string
	UserID   string
}

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *Session) IsExpired() bool {
	return s.CreatedAt.After(time.Now())
}

type Message struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	GroupID string    `json:"group_id"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}
