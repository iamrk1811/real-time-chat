package repo

import (
	"context"

	"github.com/iamrk1811/real-time-chat/types"
)

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