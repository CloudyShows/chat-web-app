// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
    server := NewServer()

    http.HandleFunc("/ws", server.handler)
    http.HandleFunc("/getUsername", server.getUsername)
    http.HandleFunc("/clear", server.clearChatHistoryHTTPHandler)

    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
        handlers.AllowedHeaders([]string{"Content-Type"}),
    )

    log.Println("Starting server on port 3001...")
    log.Fatal(http.ListenAndServe(":3001", corsHandler(http.DefaultServeMux)))
}