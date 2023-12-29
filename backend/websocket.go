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
	username, _ := s.getUsernameFromRedis(ctx, clientIP)

	s.clientData.mutex.Lock()
	s.clientData.clients[clientIP] = &ClientInfo{Conn: conn}
	s.clientData.usernames[conn] = username
	s.clientData.mutex.Unlock()

	log.Println("Client connected:", clientIP)
	s.sendChatHistory(ctx, conn)
	s.updateClientUsername(ctx, conn, r, username)
	s.broadcastUserList()
}

func (s *Server) removeClient(conn *websocket.Conn) {
	s.clientData.mutex.Lock()
	defer s.clientData.mutex.Unlock()

	delete(s.clientData.clients, conn.RemoteAddr().String())
	delete(s.clientData.usernames, conn)

	s.broadcastUserList()
}

func (s *Server) handleIncomingMessages(conn *websocket.Conn, closeChan chan struct{}) {
	defer close(closeChan) // Signal that this connection is closing

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			s.removeClient(conn)
			return
		}
		s.processMessage(conn, msg)
	}
}

func (s *Server) processMessage(conn *websocket.Conn, msg []byte) {
	var chatMessage ChatMessage
	if err := json.Unmarshal(msg, &chatMessage); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	chatMessage.Type = "message"
	chatMessage.Timestamp = time.Now()
	s.storeAndBroadcastMessage(chatMessage)
}

func (s *Server) storeAndBroadcastMessage(chatMessage ChatMessage) {
	updatedMsg, err := json.Marshal(chatMessage)
	if err != nil {
		log.Println("Error marshalling updated message:", err)
		return
	}

	// Assuming you have a context (ctx) to use with Redis operations
	ctx := context.Background() // Replace with your actual context
	if err := s.rdb.LPush(ctx, "chatHistory", string(updatedMsg)).Err(); err != nil {
		log.Println("Error saving message to Redis:", err)
	}

	s.clientData.mutex.Lock()
	defer s.clientData.mutex.Unlock()

	for _, clientInfo := range s.clientData.clients {
		if err := clientInfo.Conn.WriteMessage(websocket.TextMessage, updatedMsg); err != nil {
			log.Println("Error writing message:", err)
			delete(s.clientData.clients, clientInfo.Conn.RemoteAddr().String())
			delete(s.clientData.usernames, clientInfo.Conn)
		}
	}
}

func (s *Server) broadcastUserList() {
	s.clientData.mutex.Lock()
	defer s.clientData.mutex.Unlock()

	userList := make([]string, 0, len(s.clientData.usernames))
	for _, username := range s.clientData.usernames {
		userList = append(userList, username)
	}

	userListJSON, _ := json.Marshal(map[string]interface{}{"type": "users", "users": userList})

	for _, clientInfo := range s.clientData.clients {
		if err := clientInfo.Conn.WriteMessage(websocket.TextMessage, userListJSON); err != nil {
			log.Println("Error writing to client:", err)
			delete(s.clientData.clients, clientInfo.Conn.RemoteAddr().String())
			delete(s.clientData.usernames, clientInfo.Conn)
		}
	}
}
