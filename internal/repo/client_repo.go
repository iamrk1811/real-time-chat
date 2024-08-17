package repo

import (
	"context"

	"github.com/iamrk1811/real-time-chat/types"
)

func (c *CRUDRepo) SaveMessage(sender, receiver int, groupID int, content string) {
	ctx := context.Background()
	if receiver != 0 {
		query := `INSERT INTO messages (sender_id, receiver_id, content) VALUES($1, $2, $3)`
		c.DB.ExecContext(ctx, query, sender, receiver, content)
	} else {
		query := `INSERT INTO messages (sender_id, group_id, content) VALUES($1, $2, $3)`
		c.DB.ExecContext(ctx, query, sender, groupID, content)
	}
}

func (c *CRUDRepo) GetChats(ctx context.Context, from, to int) ([]types.Message, types.MultiError) {
	var messages []types.Message
	var mErr types.MultiError
	query := `SELECT sender_id, receiver_id, content, sent_at FROM messages WHERE sender_id=$1 AND receiver_id=$2`
	rows, err := c.DB.QueryContext(ctx, query, from, to)
	if err != nil {
		mErr.Add(err)
		return nil, mErr
	}

	for rows.Next() {
		var row types.Message

		err := rows.Scan(&row.From, &row.To, &row.Content, &row.SentAt)
		if err != nil {
			mErr.Add(err)
			continue
		}
		messages = append(messages, row)
	}
	return messages, mErr
}

func (c *CRUDRepo) GetGroupChats(ctx context.Context, sender, groupID int) ([]types.Message, types.MultiError) {
	var messages []types.Message
	var mErr types.MultiError
	query := `
	SELECT m.sender_id, m.content, m.sent_at
	FROM messages m
	JOIN UserGroups ug ON ug.group_id = m.group_id
	WHERE m.group_id = $1
	AND m.sender_id = $2
	AND EXISTS (
		SELECT 1
		FROM UserGroups ug2
		WHERE ug2.user_id = m.sender_id
		AND ug2.group_id = m.group_id
	)
	`
	rows, err := c.DB.QueryContext(ctx, query, groupID, sender)
	if err != nil {
		mErr.Add(err)
		return nil, mErr
	}

	for rows.Next() {
		var row types.Message

		err := rows.Scan(&row.From, &row.To, &row.Content, &row.SentAt)
		if err != nil {
			mErr.Add(err)
			continue
		}
		messages = append(messages, row)
	}
	return messages, mErr
}

func (c *CRUDRepo) GetUsersFromUsingGroupID(groupID int, senderUserID int) ([]types.User, error) {
	ctx := context.Background()
	var users []types.User
	query := `
	SELECT u.username
	FROM groups g
	LEFT JOIN usergroups ug ON ug.group_id = g.group_id
	LEFT JOIN users u ON u.user_id = ug.user_id
	WHERE g.group_id = $1
	AND EXISTS (
		SELECT 1
		FROM usergroups ug2
		WHERE ug2.user_id = $2
		AND ug2.group_id = g.group_id
	)
	`
	rows, err := c.DB.QueryContext(ctx, query, groupID, senderUserID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var row types.User

		err := rows.Scan(&row.Username)
		if err != nil {
			continue
		}
		users = append(users, row)
	}

	return users, nil
}
