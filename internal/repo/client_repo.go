package repo

import (
	"context"

	"github.com/iamrk1811/real-time-chat/types"
)

func (c *CRUDRepo) SaveMessage(sender, receiver string, groupID int, content string) {
	ctx := context.Background()
	query := `INSERT INTO messages (sender_id, receiver_id, group_id, content) VALUES($1, $2, $3, $4)`

	c.DB.ExecContext(ctx, query, sender, receiver, content)
}

func (c *CRUDRepo) GetChats(ctx context.Context, from, to string) ([]types.Message, types.MultiError) {
	var messages []types.Message
	var mErr types.MultiError
	query := `SELECT sender_id, receiver_id, content, sent_at FROM messages WHERE sender_id=$1 AND receiver_id=$2;`
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

func (c *CRUDRepo) GetGroupChats(group string) {

}

func (c *CRUDRepo) GetUsersFromUsingGroupID(groupID int) ([]types.User, error) {
	var users []types.User
	return users, nil
}
