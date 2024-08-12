package repo

import (
	"context"
	"time"

	"github.com/iamrk1811/real-time-chat/types"
)

func (c *CRUDRepo) GetUser(username string, password string) *types.User {
	ctx := context.TODO()
	query := `SELECT user_id, username, password FROM users WHERE username=$1 AND password=$2`
	var user types.User
	if err := c.DB.QueryRowContext(ctx, query, username, password).Scan(&user.UserID, &user.Username, &user.Password); err != nil {
		return nil
	}
	return &user
}

func (c *CRUDRepo) SaveSession(sessionID string, user types.User, expireAt time.Time) {
	ctx := context.TODO()
	query := `INSERT INTO sessions (session_id, user_id, expires_at)
              VALUES ($1, $2, $3)`

	c.DB.ExecContext(ctx, query, sessionID, user.UserID, expireAt)
}
