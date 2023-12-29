// websocket.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func (s *Server) handleNewClient(conn *websocket.Conn, r *http.Request) {
    closeChan := make(chan struct{})
    s.registerClient(conn, closeChan)

    // Start a goroutine to handle incoming messages
    go s.handleIncomingMessages(conn, closeChan)

	log.Println("Client connected:", conn.RemoteAddr())
	s.sendChatHistory(conn)

	clientIP := getClientIP(r)
	username, _ := s.getUsernameFromRedis(clientIP)
	s.updateClientUsername(clientIP, username)

	s.broadcastUserList()
}

func (s *Server) registerClient(conn *websocket.Conn, closeChan chan struct{}) {
    s.clientMutex.Lock()
    defer s.clientMutex.Unlock()

    s.clients[conn] = closeChan
}

func (s *Server) removeClient(conn *websocket.Conn) {
    s.clientMutex.Lock()
    defer s.clientMutex.Unlock()

    delete(s.clients, conn)
    delete(s.usernames, conn)

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

	if err := s.rdb.LPush(s.ctx, "chatHistory", string(updatedMsg)).Err(); err != nil {
		log.Println("Error saving message to Redis:", err)
	}

	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	for client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, updatedMsg); err != nil {
			log.Println("Error writing message:", err)
			delete(s.clients, client)
			delete(s.usernames, client)
		}
	}
}

func (s *Server) broadcastUserList() {
    s.clientMutex.Lock()
    defer s.clientMutex.Unlock()

    userList := make([]string, 0, len(s.usernames))
    for _, username := range s.usernames {
        userList = append(userList, username)
    }

    userListJSON, _ := json.Marshal(map[string]interface{}{"type": "users", "users": userList})

    for client, closeChan := range s.clients {
        select {
        case <-closeChan:
            // Connection is closed, remove it
            delete(s.clients, client)
            delete(s.usernames, client)
        default:
            if err := client.WriteMessage(websocket.TextMessage, userListJSON); err != nil {
                log.Println("Error writing to client:", err)
                delete(s.clients, client)
                delete(s.usernames, client)
            }
        }
    }
}
