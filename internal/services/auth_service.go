package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/utils"
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

	user := u.repo.GetUser(r.Context(), creds.UserName, creds.Password)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New().String()

	expireAt := time.Now().Add(3600 * time.Hour)

	go u.repo.SaveSession(context.Background(), sessionID, user, expireAt)

	cookie := http.Cookie{
		Name:     string(config.SessionKey),
		Value:    sessionID,
		Expires:  expireAt,
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)

	utils.WriteResponse(w, http.StatusOK, map[string]string{string(config.SessionKey): sessionID}, nil)
}
