package services

import (
	"encoding/json"
	"net/http"

	"github.com/iamrk1811/real-time-chat/internal/repo"
)

type Auth interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type auth struct {
	repo repo.CRUDRepo
}

func NewAuthService(repo repo.CRUDRepo) *auth {
	return &auth{
		repo: repo,
	}
}

type UserCreds struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (u *auth) Login(w http.ResponseWriter, r *http.Request) {
	var creds UserCreds

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := u.repo.GetUser(creds.UserName, creds.Password)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// sessionID := uuid.New()
	// expireAt := time.Now().Add(24 * )
	// expireAt

	// u.repo.SaveSession(sessionID, user)
}
