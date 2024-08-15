package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/services"
)

type ClientWebSocketRoutes struct {
	service services.Client
}

func NewClientWebSocketRoutes(router *mux.Router, service services.Client, config *config.Config) *ClientWebSocketRoutes {
	c := &ClientWebSocketRoutes{
		service: service,
	}
	router.HandleFunc("/chat", c.handleChat).Methods("GET")
	// config.ProtectedPaths.Add("/ws/chat")
	return c
}

func (c *ClientWebSocketRoutes) handleChat(w http.ResponseWriter, r *http.Request) {
	c.service.Chat(w, r)
}
