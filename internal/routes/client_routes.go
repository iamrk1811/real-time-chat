package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type clientRoutes struct {
	service services.Client
}

func NewClientRoutes(router *mux.Router, client services.Client) *clientRoutes {
	c := &clientRoutes{
		service: client,
	}
	router.HandleFunc("/chat", c.handleChat).Methods("GET")
	router.HandleFunc("/chat/group", c.handleGroupChat).Methods("GET")
	return c
}

func (c *clientRoutes) handleChat(w http.ResponseWriter, r *http.Request) {
	c.service.UserToUserChat(w, r)
}

func (c *clientRoutes) handleGroupChat(w http.ResponseWriter, r *http.Request) {
	c.service.GroupChat(w, r)
}
