package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	if session.IsExpired() {
		conn.Close()
		return
	}

	go c.handleConn(conn, session)
}

func (c *client) closeAndDelete(conn *websocket.Conn, key string) {
	conn.Close()
	delete(c.conns, key)
}

func (c *client) handleConn(conn *websocket.Conn, session *types.Session) {
	c.conns[session.UserName] = conn

	// setting pong handler
	conn.SetPongHandler(func(data string) error {
		return nil
	})

	// each 5 send server will send a ping
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.closeAndDelete(conn, session.UserName)
				return
			}
		}
	}()

	go c.readMessages(conn, session)
}

func (c *client) readMessages(conn *websocket.Conn, session *types.Session) {
	for {
		if session.IsExpired() {
			c.closeAndDelete(conn, session.UserName)
			return
		}

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			c.closeAndDelete(conn, session.UserName)
			break
		}

		var payload chatPayload
		if err := json.Unmarshal(p, &payload); err != nil {
			continue
		}
		go c.repo.SaveMessage(session.UserID, payload.To, payload.Content)

		// does user have active connection
		receiverConn, exist := c.conns[payload.To]
		if !exist {
			continue
		}

		if err := receiverConn.WriteMessage(messageType, p); err != nil {
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
