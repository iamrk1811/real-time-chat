package services

import (
	"net/http"

	"github.com/iamrk1811/real-time-chat/internal/repo"
)

type User interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type user struct {
	repo repo.CRUDRepo
}

func NewUserService(repo repo.CRUDRepo) *user {
	return &user{
		repo: repo,
	}
}

func (u *user) CreateUser(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("user_name")
	u.repo.CreateUser(userName)
}
