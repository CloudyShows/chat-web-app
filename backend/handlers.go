// handlers.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    conn, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
        // log.Println("Error upgrading to WebSocket:", err)
        http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    s.handleNewClient(ctx, conn, r)
}

func (s *Server) getUsername(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    clientIP := getClientIP(r)
    username, err := s.getUsernameFromRedis(ctx, clientIP)
    // log.Println("getUsername started")
    if err != nil {
        // log.Println("Error fetching username:", err)
        http.Error(w, "Failed to fetch username", http.StatusInternalServerError)
        return
    }
    // log.Println("getting username from request")

    clientSentUsername := r.URL.Query().Get("username")
    if clientSentUsername != "" {
        // log.Println("clientSentUsername")
        username = clientSentUsername
        err = s.updateUsernameInRedis(ctx, clientIP, username)
        if err != nil {
            // log.Println("Error updating username in Redis:", err)
            http.Error(w, "Failed to update username", http.StatusInternalServerError)
            return
        }
        // log.Println("Sending username to client:", username)
    }
    // log.Println(("running updateClientUsername"))

    conn, _ := s.getClientByIP(clientIP)
    s.updateClientUsername(ctx, conn, r, username)
    s.broadcastUserList()

    // log.Println("Sending username to client:", username)
    json.NewEncoder(w).Encode(map[string]string{"username": username})
}

func (s *Server) updateClientUsername(ctx context.Context, conn *websocket.Conn, r *http.Request, username string) {
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
        // log.Println("Error decoding request body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Update the username in Redis
    if err := s.updateUsernameInRedis(ctx, clientIP, reqBody.NewUsername); err != nil {
        // log.Println("Error updating username in Redis:", err)
        http.Error(w, "Failed to update username in Redis", http.StatusInternalServerError)
        return
    }

    // Retrieve the WebSocket connection associated with the client IP
    conn, err := s.getClientByIP(clientIP)
    if err != nil {
        // log.Println("Error retrieving client by IP:", err)
        http.Error(w, "Client connection not found", http.StatusInternalServerError)
        return
    }

    // Update the username in the server's clientData
    s.updateClientUsername(ctx, conn, r, reqBody.NewUsername)
    s.broadcastUserList()

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}



func (s *Server) clearChatHistoryHTTPHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.clearChatHistory(ctx)
	w.Write([]byte("Chat history cleared."))
}

func getClientIP(r *http.Request) string {
    splitIP := strings.Split(r.RemoteAddr, ":")
    return splitIP[0]
}

func (s *Server) getClientByIP(clientIP string) (*websocket.Conn, error) {
    s.clientData.mutex.Lock()
    defer s.clientData.mutex.Unlock()

    clientInfo, ok := s.clientData.clients[clientIP]
    if !ok {
        return nil, fmt.Errorf("no client found for IP: %s", clientIP)
    }

    return clientInfo.Conn, nil
}