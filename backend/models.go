// models.go
package main

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
	rdb      *redis.Client

	clientData struct {
		mutex     sync.Mutex
		clients   map[string]*ClientInfo
		usernames map[*websocket.Conn]string
	}

	chatData struct {
		mutex sync.Mutex
	}
}

type ClientInfo struct {
	Conn     *websocket.Conn
	Closed   *bool
	Username string
	CloseChan chan struct{}
}

type HeartbeatMessage struct {
    Type      string    `json:"type"`
    Timestamp time.Time `json:"timestamp"`
}

type ChatMessage struct {
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type ErrorMessage struct {
    Type      string    `json:"type"`
    Error     string    `json:"error"`
    Timestamp time.Time `json:"timestamp"`
}

type SuccessMessage struct {
    Type      string    `json:"type"`
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
}