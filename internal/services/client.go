package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/utils"
)

type Client interface {
	GetChats(w http.ResponseWriter, r *http.Request)
	GetGroupChats(w http.ResponseWriter, r *http.Request)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	repo repo.CRUDRepo
}

type chatPayload struct {
	From string `json:"from"`
	To   string `json:"to"`
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

func (c *client) GetChats(w http.ResponseWriter, r *http.Request) {
	var payload chatPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, nil, err)
		return
	}

	chats, mErr := c.repo.GetChats(r.Context(), payload.From, payload.To)
	if mErr.HasError() {
		utils.WriteResponse(w, http.StatusInternalServerError, nil, &mErr)
		return
	}

	utils.WriteResponse(w, http.StatusOK, chats, nil)
}

func (c *client) GetGroupChats(w http.ResponseWriter, r *http.Request) {

}
