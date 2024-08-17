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

// Client interface defines the methods for chat services.
type Client interface {
	GetChats(w http.ResponseWriter, r *http.Request)
	GetGroupChats(w http.ResponseWriter, r *http.Request)
	Chat(w http.ResponseWriter, r *http.Request)
}

// ClientConnectionMap is a map to store active websocket connections by user ID.
type ClientConnectionMap map[int]*websocket.Conn

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// client struct holds the dependencies required for chat services.
type client struct {
	repo   repo.CRUDRepo           // Repository interface for CRUD operations.
	conns  ClientConnectionMap      // Active WebSocket connections.
	config config.Config            // Application configuration.
}

// Payload structures for different chat functionalities.
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

// NewClientService returns a new instance of client service.
func NewClientService(repo repo.CRUDRepo, confg config.Config) *client {
	return &client{
		repo:   repo,
		conns:  ClientConnectionMap{},
		config: confg,
	}
}

// Chat handles incoming WebSocket connections and initiates the communication.
func (c *client) Chat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming connection")
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP connection to WebSocket.
	if err != nil {
		return
	}

	// Retrieve session from context and check if it is expired.
	session := r.Context().Value(config.SessionKey).(*types.Session)
	if session.IsExpired() {
		conn.Close()
		return
	}

	// Handle the connection and start communication.
	go c.handleConn(conn, session)
}

// closeAndDelete closes the WebSocket connection and removes it from the active connections map.
func (c *client) closeAndDelete(conn *websocket.Conn, key int) {
	conn.Close()
	delete(c.conns, key)
}

// handlePingPong handles ping-pong mechanism to keep the connection alive.
func (c *client) handlePingPong(conn *websocket.Conn, session *types.Session) {
	// Set pong handler to handle pong messages from the client.
	conn.SetPongHandler(func(data string) error {
		return nil
	})

	// Send ping messages every 5 seconds to keep the connection alive.
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				// Close and delete the connection if ping fails.
				c.closeAndDelete(conn, session.UserID)
				return
			}
		}
	}()
}

// handleConn handles the WebSocket connection by adding it to the active connections map
// and starts reading messages from the connection.
func (c *client) handleConn(conn *websocket.Conn, session *types.Session) {
	// Store the connection in the active connections map.
	c.conns[session.UserID] = conn

	// Start the ping-pong mechanism to keep the connection alive.
	c.handlePingPong(conn, session)

	// Start reading messages from the WebSocket connection.
	go c.readMessages(conn, session)
}

// readMessages reads messages from the WebSocket connection and handles them based on the type.
func (c *client) readMessages(conn *websocket.Conn, session *types.Session) {
	for {
		// Close the connection if the session is expired.
		if session.IsExpired() {
			c.closeAndDelete(conn, session.UserID)
			return
		}

		// Read message from the WebSocket connection.
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			c.closeAndDelete(conn, session.UserID)
			break
		}

		// Unmarshal the received message into the chat payload structure.
		var payload chatPayload
		if err := json.Unmarshal(p, &payload); err != nil {
			continue
		}

		// Save the message to the repository.
		go c.repo.SaveMessage(session.UserID, payload.To, payload.GroupID, payload.Content)

		// Handle user-to-user or user-to-group message based on the payload.
		if payload.To != 0 {
			// User-to-user message
			go c.sendToUserMessage(payload.To, messageType, p)
		} else {
			// User-to-group message
			go c.sendGroupMessage(conn, session, payload.GroupID, messageType, p)
		}
	}
}

// sendToUserMessage sends a message to a specific user if they have an active connection.
func (c *client) sendToUserMessage(to int, messageType int, p []byte) {
	// Check if the user has an active connection.
	receiverConn, exist := c.conns[to]
	if !exist {
		return
	}
	// Send the message to the user's WebSocket connection.
	receiverConn.WriteMessage(messageType, p)
}

// sendGroupMessage sends a message to all users in a specific group.
func (c *client) sendGroupMessage(senderConn *websocket.Conn, session *types.Session, groupID int, messageType int, p []byte) {
	// Retrieve all users in the group except the sender.
	groupUsers, err := c.repo.GetUsersFromUsingGroupID(groupID, session.UserID)
	if err != nil {
		senderConn.WriteMessage(websocket.TextMessage, []byte(config.MessageFailed))
		return
	}
	// Send the message to all users in the group who have an active connection.
	for _, u := range groupUsers {
		receiverConn, exist := c.conns[u.UserID]
		if !exist {
			continue
		}
		receiverConn.WriteMessage(messageType, p)
	}
}

// GetChats retrieves the chat history between two users.
func (c *client) GetChats(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(config.SessionKey).(*types.Session)
	if session.IsExpired() {
		utils.WriteResponse(w, http.StatusBadRequest, nil, errors.New("session expired"))
		return
	}

	// Decode the request payload to get the user IDs.
	var payload getChatsPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, nil, err)
		return
	}

	// Ensure the request is from the authenticated user.
	if payload.From != session.UserID {
		utils.WriteResponse(w, http.StatusBadRequest, nil, errors.New("resource doesn't belong to you"))
		return
	}

	// Retrieve chat history from the repository.
	chats, mErr := c.repo.GetChats(r.Context(), payload.From, payload.To)
	if mErr.HasError() {
		utils.WriteResponse(w, http.StatusInternalServerError, nil, &mErr)
		return
	}

	// Write the chat history as the response.
	utils.WriteResponse(w, http.StatusOK, chats, nil)
}

// GetGroupChats retrieves the chat history of a specific group.
func (c *client) GetGroupChats(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(config.SessionKey).(*types.Session)
	if session.IsExpired() {
		utils.WriteResponse(w, http.StatusBadRequest, nil, errors.New("session expired"))
		return
	}

	// Decode the request payload to get the group ID.
	var payload getGroupChatsPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, nil, err)
		return
	}

	// Retrieve group chat history from the repository.
	chats, mErr := c.repo.GetGroupChats(r.Context(), session.UserID, payload.GroupID)
	if mErr.HasError() {
		utils.WriteResponse(w, http.StatusInternalServerError, nil, &mErr)
		return
	}

	// Write the group chat history as the response.
	utils.WriteResponse(w, http.StatusOK, chats, nil)
}
