package repo

import (
	"context"
	"log"
)

func (u *CRUDRepo) CreateUser(userName string) {
	ctx := context.TODO()
	query := `INSERT INTO users (user_name) VALUES ($1)`
	if _, err := u.DB.ExecContext(ctx, query, userName); err != nil {
		log.Printf("failed to create user %s", userName)
	}
}
