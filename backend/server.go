// server.go
package main

import (
	"context"
	"net/http"

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
            Addr: "localhost:6379", // Redis server address
        }),
        ctx: context.Background(),
        clients: make(map[*websocket.Conn]bool),
        usernames: make(map[*websocket.Conn]string),
    }
}