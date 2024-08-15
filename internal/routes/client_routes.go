package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type clientRoutes struct {
	service services.Client
}

func NewClientRoutes(router *mux.Router, client services.Client, config *config.Config) *clientRoutes {
	c := &clientRoutes{
		service: client,
	}

	router.HandleFunc("/user/chats", c.handleChat)
	config.ProtectedPaths.Add("/api/user/chats")

	router.HandleFunc("/user/chats/group", c.handleGroupChat)
	config.ProtectedPaths.Add("/api/user/chats/group")
	return c
}

func (c *clientRoutes) handleChat(w http.ResponseWriter, r *http.Request) {
	c.service.GetChats(w, r)
}

func (c *clientRoutes) handleGroupChat(w http.ResponseWriter, r *http.Request) {
	c.service.GetGroupChats(w, r)
}
