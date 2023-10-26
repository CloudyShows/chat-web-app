// handlers.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
	ctx = context.Background()
)

var mutex = &sync.Mutex{}                        // Variable for synchronization
var clients = make(map[*websocket.Conn]bool)     // connected clients
var usernames = make(map[*websocket.Conn]string) // track usernames

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	defer conn.Close()

	// Lock the mutex before updating maps (New Code)
	mutex.Lock()
	clients[conn] = true

	log.Println("Client connected:", conn.RemoteAddr())

	// Fetch chat history from Redis when a new client connects
	chatHistory, err := rdb.LRange(ctx, "chatHistory", 0, -1).Result()
	if err != nil {
		log.Println("Error fetching chat history:", err)
	} else {
		for _, msg := range chatHistory {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Println("Error sending chat history:", err)
				return
			}
		}
	}

	splitIP := strings.Split(r.RemoteAddr, ":")
	clientIP := splitIP[0]

	log.Println("Client IP:", clientIP) // Debug line

	username, err := rdb.Get(ctx, "username:"+clientIP).Result()
	log.Println("Fetched username from Redis:", username) // Debug line

	if err == redis.Nil {
		username = "User"
	} else if err != nil {
		log.Println("Error fetching username:", err)
	}

	clientSentUsername := r.URL.Query().Get("username")
	log.Println("Client sent username:", clientSentUsername) // Debug line

	if clientSentUsername != "" {
		username = clientSentUsername
		log.Println("Client sent username:", username) // Debug line
		log.Println("Client IP:", clientIP)           // Debug line
		err = updateUsernameInRedis(clientIP, username)
		if err != nil {
			log.Println("Error updating username in Redis:", err)
		}
		log.Println("Updated username in Redis:", username) // Debug line
	}

	// Update the usernames map
	usernames[conn] = username
	// Unlock the mutex after updating (New Code)
	mutex.Unlock()
	// Broadcast the list of connected users
	broadcastUserList()

	// Remove client and username when they disconnect
	defer func() {
		delete(clients, conn)
		delete(usernames, conn)
		broadcastUserList()
	}()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var chatMessage ChatMessage
		if err := json.Unmarshal(msg, &chatMessage); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		// Add the message type
		chatMessage.Type = "message"

		// Re-encode the message with the added type field
		updatedMsg, err := json.Marshal(chatMessage)
		if err != nil {
			log.Println("Error marshalling updated message:", err)
			continue
		}

		// Save message to Redis
		err = rdb.LPush(ctx, "chatHistory", string(updatedMsg)).Err()
		if err != nil {
			log.Println("Error saving message to Redis:", err)
		}

		// Broadcast message to all connected clients
		for client := range clients {
			if err := client.WriteMessage(messageType, updatedMsg); err != nil {
				log.Println("Error writing message:", err)
				delete(clients, client)
				delete(usernames, client)
			}
		}

	}
}

// Function to broadcast the list of connected users
func broadcastUserList() {
	// Lock the mutex before reading from the map (New Code)
	mutex.Lock()
	userList := make([]string, 0, len(usernames))
	for _, username := range usernames {
		userList = append(userList, username)
	}
	// Unlock the mutex after reading (New Code)
	mutex.Unlock()
	userListJSON, _ := json.Marshal(map[string]interface{}{"type": "users", "users": userList})
	for client := range clients {
		client.WriteMessage(websocket.TextMessage, userListJSON)
	}
}

func getUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	splitIP := strings.Split(r.RemoteAddr, ":")
	clientIP := splitIP[0]
	username, err := rdb.Get(ctx, "username:"+clientIP).Result()

	if err == redis.Nil {
		username = "User"
	} else if err != nil {
		log.Println("Error fetching username:", err)
	}

	clientSentUsername := r.URL.Query().Get("username")
	if clientSentUsername != "" {
		username = clientSentUsername
		err = updateUsernameInRedis(clientIP, username)
		if err != nil {
			log.Println("Error updating username in Redis:", err)
		}
	}

	// Update the usernames map for the corresponding WebSocket connection
	for conn := range clients {
		if conn.RemoteAddr().String() == clientIP {
			usernames[conn] = username
			break
		}
	}

	// Broadcast the updated list of usernames
	broadcastUserList()

	json.NewEncoder(w).Encode(map[string]string{"username": username})
}

func clearChatHistoryHTTPHandler(w http.ResponseWriter, r *http.Request) {
	clearChatHistory()
	w.Write([]byte("Chat history cleared."))
}
