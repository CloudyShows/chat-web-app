// redis.go
package main

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

func (s *Server) sendChatHistory(conn *websocket.Conn) {
	// Retrieve chat history from Redis
	chatHistory, err := s.rdb.LRange(s.ctx, "chatHistory", 0, -1).Result()
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
func (s *Server) getUsernameFromRedis(clientIP string) (string, error) {
    username, err := s.rdb.Get(s.ctx, "username:"+clientIP).Result()
    if err == redis.Nil {
        return "User", nil
    }
    return username, err
}

func (s *Server) updateUsernameInRedis(clientIP string, username string) error {
	return s.rdb.Set(s.ctx, "username:"+clientIP, username, 0).Err()
}

func (s *Server) clearChatHistory() {
	err := s.rdb.Del(s.ctx, "chatHistory").Err()
	if err != nil {
		log.Println("Error clearing chat history:", err)
	}
}