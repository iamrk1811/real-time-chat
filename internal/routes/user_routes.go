package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type userRoute struct {
	service services.User
}

func NewUserRoutes(router *mux.Router, service services.User) *userRoute {
	u := &userRoute{
		service: service,
	}
	router.HandleFunc("/user/create", u.handleCreateUser).Methods("POST")
	return u
}

func (route *userRoute) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	route.service.CreateUser(w, r)
}
