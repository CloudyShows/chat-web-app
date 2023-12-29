// server.go
package main

import (
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		rdb: redis.NewClient(&redis.Options{
			Addr: "redis:6379", // Redis container address
		}),
		clientData: struct {
			mutex     sync.Mutex
			clients   map[string]*ClientInfo
			usernames map[*websocket.Conn]string
		}{
			clients:   make(map[string]*ClientInfo),
			usernames: make(map[*websocket.Conn]string),
		},
	}
}
