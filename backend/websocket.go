// websocket.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func (s *Server) handleNewClient(ctx context.Context, conn *websocket.Conn, r *http.Request) {
	clientIP := getClientIP(r)
	username, err := s.getUsernameFromRedis(ctx, clientIP)
	if err != nil {
		log.Printf("Error retrieving username from Redis for IP %s: %v\n", clientIP, err)
		username = "anonymous" // Default to anonymous if username not found
	}

	s.clientData.mutex.Lock()
	s.clientData.clients[clientIP] = &ClientInfo{
		Conn:      conn,
		Closed:    false, // Set the Closed field to false directly though it's the zero value by default
		CloseChan: make(chan struct{}),
		Username:  username,
	}
	s.clientData.usernames[conn] = username
	s.clientData.mutex.Unlock()

	log.Printf("Client connected: %s, Username: %s\n", clientIP, username)
	s.sendChatHistory(ctx, conn)
	s.updateClientUsername(r, username)
	s.broadcastUserList(true)
	// Start handling incoming messages
	closeChan := make(chan struct{})
	// Store the closeChan in the client's ClientInfo
	s.clientData.clients[clientIP].CloseChan = closeChan
	go s.handleIncomingMessages(conn, closeChan)
	s.sendSuccessMessage(conn, "Connected to chat server")
}

func (s *Server) removeClient(conn *websocket.Conn, lock bool) {
	if lock {
		s.clientData.mutex.Lock()
		defer s.clientData.mutex.Unlock()
	}
	clientIP := getClientIP(conn)
	clientInfo, ok := s.clientData.clients[clientIP]
	if !ok {
		log.Printf("Client %s not found\n", clientIP)
		return
	}
	// Signal to stop handling incoming messages
	close(clientInfo.CloseChan)
	clientInfo.Closed = true // Set the Closed field to true directly
	delete(s.clientData.clients, clientIP)
	delete(s.clientData.usernames, conn)
	log.Printf("Client disconnected: %s\n", clientIP)
	if lock {
		s.broadcastUserList(false)
	}
}

func (s *Server) handleIncomingMessages(conn *websocket.Conn, closeChan chan struct{}) {
	defer s.removeClient(conn, true) // Ensure the lock is used when removing the client

	for {
		select {
		case <-closeChan:
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading message: %v\n", err)
				}
				return
			}
			go s.processMessage(conn, msg) // Process each message in a separate goroutine
		}
	}
}

func (s *Server) processMessage(conn *websocket.Conn, msg []byte) {
	// log.Printf("Message received: %s\n", string(msg))
	var chatMessage ChatMessage
	if err := json.Unmarshal(msg, &chatMessage); err != nil {
		log.Printf("Error unmarshalling JSON: %v\n", err)
		s.sendErrorMessage(conn, "Error unmarshalling JSON")
		return
	}

	switch chatMessage.Type {
	case "message":
		chatMessage.Timestamp = time.Now()
		log.Printf("Message received: %s\n", string(msg))
		if err := s.storeAndBroadcastMessage(conn, chatMessage); err != nil {
			s.sendErrorMessage(conn, "Error storing and broadcasting message")
		}
	// Add cases for other message types here
	default:
		log.Printf("Unknown message type: %s\n", chatMessage.Type)
	}
}

func (s *Server) storeAndBroadcastMessage(conn *websocket.Conn, chatMessage ChatMessage) error {
	updatedMsg, err := json.Marshal(chatMessage)
	if err != nil {
		log.Printf("Error marshalling message: %v\n", err)
		return err
	}

	ctx := context.Background() // Consider using a more specific context if available
	if err := s.rdb.LPush(ctx, "chatHistory", string(updatedMsg)).Err(); err != nil {
		log.Printf("Error saving message to Redis: %v\n", err)
		return err
	}

	log.Printf("Message stored: %s\n", string(updatedMsg))
	s.broadcastMessage(updatedMsg)
	return nil
}

func (s *Server) broadcastMessage(message []byte) {
	s.clientData.mutex.Lock()
	defer s.clientData.mutex.Unlock()

	for clientIP, clientInfo := range s.clientData.clients {
		if err := clientInfo.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message to client %s: %v\n", clientIP, err)
			s.removeClient(clientInfo.Conn, false)
		}
	}
}

func (s *Server) broadcastUserList(lock bool) {
	if lock {
		s.clientData.mutex.Lock()
		defer s.clientData.mutex.Unlock()
	}

	userListJSON, err := s.prepareUserListJSON()
	if err != nil {
		log.Printf("Error preparing user list JSON: %v\n", err)
		return
	}

	for clientIP, clientInfo := range s.clientData.clients {
		if clientInfo.Conn == nil || clientInfo.Closed {
			continue
		}
		if err := clientInfo.Conn.WriteMessage(websocket.TextMessage, userListJSON); err != nil {
			log.Printf("Error writing user list to client %s: %v\n", clientIP, err)
			s.removeClient(clientInfo.Conn, false) // This should now set Closed = true
		}
	}
}

func (s *Server) prepareUserListJSON() ([]byte, error) {
	userList := make([]string, 0, len(s.clientData.usernames))
	for _, username := range s.clientData.usernames {
		userList = append(userList, username)
	}
	return json.Marshal(map[string]interface{}{"type": "users", "users": userList})
}

func (s *Server) sendErrorMessage(conn *websocket.Conn, message string) {
	errorMessage := ErrorMessage{
		Type:      "error",
		Error:     message,
		Timestamp: time.Now(),
	}
	errorMsg, err := json.Marshal(errorMessage)
	if err != nil {
		log.Printf("Error marshalling error message: %v\n", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, errorMsg); err != nil {
		log.Printf("Error sending error message to client: %v\n", err)
	}
}

func (s *Server) sendSuccessMessage(conn *websocket.Conn, message string) {
	successMessage := SuccessMessage{
		Type:      "success",
		Message:   message,
		Timestamp: time.Now(),
	}
	successMsg, err := json.Marshal(successMessage)
	if err != nil {
		log.Printf("Error marshalling success message: %v\n", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, successMsg); err != nil {
		log.Printf("Error sending success message to client: %v\n", err)
	}
}
