package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func (s *Server) handleNewClient(conn *websocket.Conn, r *http.Request) {
    s.mutex.Lock()
    s.clients[conn] = true
    s.mutex.Unlock()

    log.Println("Client connected:", conn.RemoteAddr())
    s.sendChatHistory(conn)
    s.updateUsername(conn, r)
    s.broadcastUserList()
}

func (s *Server) handleIncomingMessages(conn *websocket.Conn) {
    for {
        messageType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            s.mutex.Lock()
            delete(s.clients, conn)
            delete(s.usernames, conn)
            s.mutex.Unlock()
            return
        }
        s.processMessage(conn, messageType, msg)
    }
}

func (s *Server) processMessage(conn *websocket.Conn, messageType int, msg []byte) {
    var chatMessage ChatMessage
    if err := json.Unmarshal(msg, &chatMessage); err != nil {
        log.Println("Error unmarshalling JSON:", err)
        return
    }

    chatMessage.Type = "message"
    chatMessage.Timestamp = time.Now()

    updatedMsg, err := json.Marshal(chatMessage)
    if err != nil {
        log.Println("Error marshalling updated message:", err)
        return
    }

    err = s.rdb.LPush(s.ctx, "chatHistory", string(updatedMsg)).Err()
    if err != nil {
        log.Println("Error saving message to Redis:", err)
    }

    s.clientMutex.Lock()
    for client := range s.clients {
        if err := client.WriteMessage(messageType, updatedMsg); err != nil {
            log.Println("Error writing message:", err)
            delete(s.clients, client)
            delete(s.usernames, client)
        }
    }
    s.clientMutex.Unlock()
}

func (s *Server) broadcastUserList() {
    s.mutex.Lock()
    userList := make([]string, 0, len(s.usernames))
    for _, username := range s.usernames {
        userList = append(userList, username)
    }
    s.mutex.Unlock()

    userListJSON, _ := json.Marshal(map[string]interface{}{"type": "users", "users": userList})

    s.clientMutex.Lock()
    for client := range s.clients {
        if err := client.WriteMessage(websocket.TextMessage, userListJSON); err != nil {
            log.Println("Error writing to client:", err)
        }
    }
    s.clientMutex.Unlock()
}