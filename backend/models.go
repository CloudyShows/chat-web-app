// models.go
package main

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Server struct {
    upgrader websocket.Upgrader
    rdb      *redis.Client
    ctx      context.Context
    clients  map[*websocket.Conn]bool
    usernames map[*websocket.Conn]string
    mutex sync.Mutex
    clientMutex sync.Mutex
}

type ChatMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Timestamp time.Time `json:"timestamp"`

}