package repo

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"github.com/iamrk1811/real-time-chat/config"
)

type CRUDRepo struct {
	*sql.DB
}

func NewCRUDRepo(config config.Config) *CRUDRepo {
	db, err := sql.Open("postgres", config.DB.CONN_URL)
	if err != nil {
		log.Fatalf("failed to connect db %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("db is not ready yet %v", err)
	}

	return &CRUDRepo{db}
}
