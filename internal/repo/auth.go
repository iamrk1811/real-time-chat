package repo

import (
	"context"
	"time"

	"github.com/iamrk1811/real-time-chat/types"
)

func (c *CRUDRepo) GetUser(ctx context.Context, username string, password string) *types.User {
	query := `SELECT user_id, username, password FROM users WHERE username=$1 AND password=$2`
	var user types.User
	if err := c.DB.QueryRowContext(ctx, query, username, password).Scan(&user.UserID, &user.Username, &user.Password); err != nil {
		return nil
	}
	return &user
}

func (c *CRUDRepo) SaveSession(ctx context.Context, sessionID string, user *types.User, expireAt time.Time) {
	query := `INSERT INTO sessions (session_id, user_id, expires_at)
              VALUES ($1, $2, $3)`

	c.DB.ExecContext(ctx, query, sessionID, user.UserID, expireAt)
}

func (c *CRUDRepo) FetchUserBySessionID(ctx context.Context, sessionID string) *types.Session {
	query := `SELECT u.user_id, session_id, s.created_at, s.expires_at FROM sessions AS s LEFT JOIN users AS u ON s.user_id=u.user_id WHERE s.session_id=$1;`
	var session types.Session
	if err := c.DB.QueryRowContext(ctx, query, sessionID).Scan(&session.UserID, &session.SessionID, &session.CreatedAt, &session.ExpiresAt); err != nil {
		return nil
	}
	return &session
}
