package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
    conn, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading to WebSocket:", err)
        http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    s.handleNewClient(conn, r)
}

func (s *Server) getUsername(w http.ResponseWriter, r *http.Request) {
    clientIP := getClientIP(r)
    username, err := s.getUsernameFromRedis(clientIP)
    if err != nil {
        log.Println("Error fetching username:", err)
        http.Error(w, "Failed to fetch username", http.StatusInternalServerError)
        return
    }

    clientSentUsername := r.URL.Query().Get("username")
    if clientSentUsername != "" {
        username = clientSentUsername
        err = s.updateUsernameInRedis(clientIP, username)
        if err != nil {
            log.Println("Error updating username in Redis:", err)
            http.Error(w, "Failed to update username", http.StatusInternalServerError)
            return
        }
    }

    s.updateClientUsername(clientIP, username)
    s.broadcastUserList()

    json.NewEncoder(w).Encode(map[string]string{"username": username})
}

func getClientIP(r *http.Request) string {
    splitIP := strings.Split(r.RemoteAddr, ":")
    return splitIP[0]
}

func (s *Server) updateClientUsername(clientIP, username string) {
    s.clientMutex.Lock()
    defer s.clientMutex.Unlock()

    for conn := range s.clients {
        if conn.RemoteAddr().String() == clientIP {
            s.usernames[conn] = username
            break
        }
    }
}

func (s *Server) clearChatHistoryHTTPHandler(w http.ResponseWriter, r *http.Request) {
    s.clearChatHistory()
    w.Write([]byte("Chat history cleared."))
}