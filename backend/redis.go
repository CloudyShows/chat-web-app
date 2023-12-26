package main

import (
	"log"

	"github.com/gorilla/websocket"
)


func (s *Server) sendChatHistory(conn *websocket.Conn) {
    chatHistory, err := s.rdb.LRange(s.ctx, "chatHistory", 0, -1).Result()
    if err != nil {
        log.Println("Error fetching chat history:", err)
        return
    }
    for _, msg := range chatHistory {
        if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
            log.Println("Error sending chat history:", err)
            return
        }
    }
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