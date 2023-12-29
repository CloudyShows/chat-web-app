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
}

type ChatMessage struct {
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}
