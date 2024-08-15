package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type authRoutes struct {
	service services.Auth
}

func NewAuthRoutes(router *mux.Router, service services.Auth) *authRoutes {
	u := &authRoutes{
		service: service,
	}
	router.HandleFunc("/auth/login", u.handleLogin).Methods("POST")
	return u
}

func (route *authRoutes) handleLogin(w http.ResponseWriter, r *http.Request) {
	route.service.Login(w, r)

}
