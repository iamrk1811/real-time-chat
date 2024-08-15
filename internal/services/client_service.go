package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iamrk1811/real-time-chat/config"
	"github.com/iamrk1811/real-time-chat/internal/repo"
	"github.com/iamrk1811/real-time-chat/types"
	"github.com/iamrk1811/real-time-chat/utils"
)

type Client interface {
	GetChats(w http.ResponseWriter, r *http.Request)
	GetGroupChats(w http.ResponseWriter, r *http.Request)
	Chat(w http.ResponseWriter, r *http.Request)
}

type ClientConnectionMap map[string]*websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	repo   repo.CRUDRepo
	conns  ClientConnectionMap
	config config.Config
}

type getChatsPayload struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type chatPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

func NewClientService(repo repo.CRUDRepo, confg config.Config) *client {
	return &client{
		repo:   repo,
		conns:  ClientConnectionMap{},
		config: confg,
	}
}

func (c *client) Chat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incomming connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	session := r.Context().Value(config.SessionKey).(*types.Session)

	go c.handleConn(conn, session)
}

func (c *client) handleConn(conn *websocket.Conn, session *types.Session) {
	c.conns[session.UserID] = conn
	go c.readMessages(conn)
}

func (c *client) readMessages(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			continue
		}
		fmt.Println("MSG", string(p))
		var payload chatPayload
		if err := json.Unmarshal(p, &payload); err != nil {
			continue
		}

		if err := c.conns[payload.To].WriteMessage(messageType, p); err != nil {
			continue
		}
	}
}

func (c *client) GetChats(w http.ResponseWriter, r *http.Request) {
	var payload getChatsPayload
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
