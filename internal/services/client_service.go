package services

import (
	"encoding/json"
	"errors"
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

type ClientConnectionMap map[int]*websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	repo repo.CRUDRepo
	// contains active websocket coonection
	conns  ClientConnectionMap
	config config.Config
}

type getChatsPayload struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type getGroupChatsPayload struct {
	GroupID int `json:"group_id"`
}

type chatPayload struct {
	To      int    `json:"to"`
	From    int    `json:"from"`
	Content string `json:"content"`
	GroupID int    `json:"group_id"`
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

func (c *client) closeAndDelete(conn *websocket.Conn, key int) {
	conn.Close()
	delete(c.conns, key)
}

func (c *client) handlePingPong(conn *websocket.Conn, session *types.Session) {
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
				c.closeAndDelete(conn, session.UserID)
				return
			}
		}
	}()
}

func (c *client) handleConn(conn *websocket.Conn, session *types.Session) {
	c.conns[session.UserID] = conn

	c.handlePingPong(conn, session)

	go c.readMessages(conn, session)
}

func (c *client) readMessages(conn *websocket.Conn, session *types.Session) {
	for {
		if session.IsExpired() {
			c.closeAndDelete(conn, session.UserID)
			return
		}

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			c.closeAndDelete(conn, session.UserID)
			break
		}

		var payload chatPayload
		if err := json.Unmarshal(p, &payload); err != nil {
			continue
		}
		go c.repo.SaveMessage(session.UserID, payload.To, payload.GroupID, payload.Content)

		if payload.To != 0 {
			// user to user message
			go c.sendToUserMessage(payload.To, messageType, p)
		} else {
			// user to group
			go c.sendGroupMessage(conn, session, payload.GroupID, messageType, p)
		}
	}
}

func (c *client) sendToUserMessage(to int, messageType int, p []byte) {
	// does user have active connection
	receiverConn, exist := c.conns[to]
	if !exist {
		return
	}
	receiverConn.WriteMessage(messageType, p)
}

func (c *client) sendGroupMessage(senderConn *websocket.Conn, session *types.Session, groupID int, messageType int, p []byte) {
	groupUsers, err := c.repo.GetUsersFromUsingGroupID(groupID, session.UserID)
	if err != nil {
		senderConn.WriteMessage(websocket.TextMessage, []byte(config.MessageFailed))
		return
	}
	for _, u := range groupUsers {
		receiverConn, exist := c.conns[u.UserID]
		if !exist {
			continue
		}
		receiverConn.WriteMessage(messageType, p)
	}
}

func (c *client) GetChats(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(config.SessionKey).(*types.Session)
	if session.IsExpired() {
		utils.WriteResponse(w, http.StatusBadRequest, nil, errors.New("session expired"))
		return
	}
	var payload getChatsPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, nil, err)
		return
	}
	// check if from and session user is the same
	if payload.From != session.UserID {
		utils.WriteResponse(w, http.StatusBadRequest, nil, errors.New("resource doesn't belogs to you"))
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
	session := r.Context().Value(config.SessionKey).(*types.Session)
	if session.IsExpired() {
		utils.WriteResponse(w, http.StatusBadRequest, nil, errors.New("session expired"))
		return
	}

	var payload getGroupChatsPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, nil, err)
		return
	}

	chats, mErr := c.repo.GetGroupChats(r.Context(), session.UserID, payload.GroupID)
	if mErr.HasError() {
		utils.WriteResponse(w, http.StatusInternalServerError, nil, &mErr)
		return
	}

	utils.WriteResponse(w, http.StatusOK, chats, nil)
}
