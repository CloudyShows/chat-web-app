// redis.go
package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

func (s *Server) sendChatHistory(ctx context.Context, conn *websocket.Conn) {

	// Retrieve chat history from Redis
	chatHistory, err := s.rdb.LRange(ctx, "chatHistory", 0, -1).Result()
	if err != nil {
		log.Println("Error fetching chat history:", err)
		return
	}

	// Reverse the chat history to send the oldest messages first
	for i, j := 0, len(chatHistory)-1; i < j; i, j = i+1, j-1 {
		chatHistory[i], chatHistory[j] = chatHistory[j], chatHistory[i]
	}

	// Send the reversed chat history to the client
	for _, msg := range chatHistory {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("Error sending chat history:", err)
			return
		}
	}
}

func (s *Server) getUsernameFromRedis(ctx context.Context, clientIP string) (string, error) {

	username, err := s.rdb.Get(ctx, "username:"+clientIP).Result()
	if err == redis.Nil {
		return "User", nil // Default username if not found
	}
	if err != nil {
		log.Printf("Error getting username from Redis for IP %s: %s", clientIP, err)
		return "", err
	}

	return username, nil
}

func (s *Server) updateUsernameInRedis(ctx context.Context, clientIP string, username string) error {

	err := s.rdb.Set(ctx, "username:"+clientIP, username, 0).Err()
	if err != nil {
		log.Printf("Error updating username in Redis for IP %s: %s", clientIP, err)
	}
	return err
}

func (s *Server) clearChatHistory(ctx context.Context) {
    err := s.rdb.Del(ctx, "chatHistory").Err()
    if err != nil {
        log.Println("Error clearing chat history:", err)
    }
}