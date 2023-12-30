// handlers.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	s.handleNewClient(ctx, conn, r)
}

func (s *Server) getUsername(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    log.Println("/getUsername API called")
	clientIP := getClientIP(r)
	username, err := s.getUsernameFromRedis(ctx, clientIP)
	if err != nil {
		log.Println("Error fetching username:", err)
		http.Error(w, "Failed to fetch username", http.StatusInternalServerError)
		return
	}

	clientSentUsername := r.URL.Query().Get("username")
	if clientSentUsername != "" {
		username = clientSentUsername
		err = s.updateUsernameInRedis(ctx, clientIP, username)
		if err != nil {
			log.Println("Error updating username in Redis:", err)
			http.Error(w, "Failed to update username", http.StatusInternalServerError)
			return
		}
	}

	s.updateClientUsername(r, username)
	s.broadcastUserList()

	log.Println("Sending username to client:", username)
	json.NewEncoder(w).Encode(map[string]string{"username": username})
}

func (s *Server) updateClientUsername(r *http.Request, username string) {
	s.clientData.mutex.Lock()
	defer s.clientData.mutex.Unlock()

	clientIP := getClientIP(r)

	if client, ok := s.clientData.clients[clientIP]; ok {
		client.Username = username
	} else {
		log.Printf("No client found for connection: %s", clientIP)
	}

	log.Printf("Username %s set for connection: %s", username, clientIP)
}

func (s *Server) changeUsernameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	clientIP := getClientIP(r)

	var reqBody struct {
		OldUsername string `json:"oldUsername"`
		NewUsername string `json:"newUsername"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the username in Redis
	if err := s.updateUsernameInRedis(ctx, clientIP, reqBody.NewUsername); err != nil {
		log.Println("Error updating username in Redis:", err)
		http.Error(w, "Failed to update username in Redis", http.StatusInternalServerError)
		return
	}

	// Update the username in the server's clientData
	s.updateClientUsername(r, reqBody.NewUsername)
	s.broadcastUserList()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}

func (s *Server) clearChatHistoryHTTPHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.clearChatHistory(ctx)
	w.Write([]byte("Chat history cleared."))
}

func getClientIP(i interface{}) string {
	var addr string

	switch v := i.(type) {
	case *http.Request:
		addr = v.RemoteAddr
	case *websocket.Conn:
		addr = v.RemoteAddr().String()
	default:
		return ""
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		// If splitting fails, return the whole string
		return addr
	}

	return host
}