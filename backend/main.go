package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

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
		// Addr: "redis:6379", // Redis server address
	})
	ctx = context.Background()
)

var clients = make(map[*websocket.Conn]bool) // connected clients

func getUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientIP := r.RemoteAddr
	username, err := rdb.Get(ctx, "username:"+clientIP).Result()

	if err == redis.Nil {
		username = "user"
	} else if err != nil {
		log.Println("Error fetching username:", err)
	}

	// Check if the client sent a username as a query parameter
	clientSentUsername := r.URL.Query().Get("username")
	if clientSentUsername != "" {
		username = clientSentUsername
		// Save it to Redis
		err = rdb.Set(ctx, "username:"+clientIP, username, 0).Err()
		if err != nil {
			log.Println("Error saving username:", err)
		}
	}

	json.NewEncoder(w).Encode(map[string]string{"username": username})
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()
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

	clientIP := r.RemoteAddr
	username, err := rdb.Get(ctx, "username:"+clientIP).Result()
	if err == redis.Nil {
		username = "defaultUsername" // or prompt the user to enter a username
	} else if err != nil {
		log.Println("Error fetching username:", err)
	}

	// Save the username to Redis
	err = rdb.Set(ctx, "username:"+clientIP, username, 0).Err()
	if err != nil {
		log.Println("Error saving username:", err)
	}

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, conn)
			return
		}

		var chatMessage ChatMessage
		if err := json.Unmarshal(msg, &chatMessage); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}
		log.Println("Received message with username:", chatMessage.Username) // Debug line

		// Save message to Redis
		err = rdb.LPush(ctx, "chatHistory", string(msg)).Err()
		if err != nil {
			log.Println("Error saving message to Redis:", err)
		}

		// Broadcast message to all connected clients
		for client := range clients {
			if err := client.WriteMessage(messageType, msg); err != nil {
				log.Println("Error writing message:", err)
				delete(clients, client)
			}
		}
	}
}

func clearChatHistory() {
	err := rdb.Del(ctx, "chatHistory").Err()
	if err != nil {
		log.Println("Error clearing chat history:", err)
	} else {
		log.Println("Chat history cleared.")
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	http.HandleFunc("/getUsername", getUsername)
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		clearChatHistory()
		w.Write([]byte("Chat history cleared."))
	})

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	log.Println("Starting server on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", corsHandler(http.DefaultServeMux)))
}
