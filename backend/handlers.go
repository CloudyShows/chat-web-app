package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)


func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
    conn, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading to WebSocket:", err)
        return
    }
    defer conn.Close()

    s.handleNewClient(conn, r)
    s.handleIncomingMessages(conn)
}

func (s *Server) getUsername(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    splitIP := strings.Split(r.RemoteAddr, ":")
    clientIP := splitIP[0]
    username, err := s.rdb.Get(s.ctx, "username:"+clientIP).Result()

    if err == redis.Nil {
        username = "User"
    } else if err != nil {
        log.Println("Error fetching username:", err)
    }

    clientSentUsername := r.URL.Query().Get("username")
    if clientSentUsername != "" {
        username = clientSentUsername
        err = s.updateUsernameInRedis(clientIP, username)
        if err != nil {
            log.Println("Error updating username in Redis:", err)
        }
    }

    s.clientMutex.Lock()
    for conn := range s.clients {
        if conn.RemoteAddr().String() == clientIP {
            s.usernames[conn] = username
            break
        }
    }
    s.clientMutex.Unlock()

    s.broadcastUserList()

    json.NewEncoder(w).Encode(map[string]string{"username": username})
}

func (s *Server) updateUsername(conn *websocket.Conn, r *http.Request) {
    splitIP := strings.Split(r.RemoteAddr, ":")
    clientIP := splitIP[0]
    username, err := s.rdb.Get(s.ctx, "username:"+clientIP).Result()

    if err == redis.Nil {
        username = "User"
    } else if err != nil {
        log.Println("Error fetching username:", err)
    }

    clientSentUsername := r.URL.Query().Get("username")
    if clientSentUsername != "" {
        username = clientSentUsername
        err = s.updateUsernameInRedis(clientIP, username)
        if err != nil {
            log.Println("Error updating username in Redis:", err)
        }
    }

    s.clientMutex.Lock()
    s.usernames[conn] = username
    s.clientMutex.Unlock()
}



func (s *Server) clearChatHistoryHTTPHandler(w http.ResponseWriter, r *http.Request) {
    s.clearChatHistory()
    w.Write([]byte("Chat history cleared."))
}