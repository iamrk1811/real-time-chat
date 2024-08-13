package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iamrk1811/real-time-chat/internal/repo"
)

type Client interface {
	UserToUserChat(w http.ResponseWriter, r *http.Request)
	GroupChat(w http.ResponseWriter, r *http.Request)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

type client struct {
	repo repo.CRUDRepo
}

func NewClientService(repo repo.CRUDRepo) *client {
	return &client{
		repo: repo,
	}
}

func (c *client) UserToUserChat(w http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("session_id")

	// if err != nil {
	// 	fmt.Println("error")
	// 	return
	// }

	fmt.Println(r.Cookies(), "Got it")
	// return
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrading fail", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading:", err)
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println("Error writing:", err)
		}
	}

}

func (c *client) GroupChat(w http.ResponseWriter, r *http.Request) {

}
