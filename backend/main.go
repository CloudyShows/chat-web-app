// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	server := NewServer()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.handler(w, r) // This will pass the context internally
	})
	http.HandleFunc("/getUsername", func(w http.ResponseWriter, r *http.Request) {
		server.getUsername(r.Context(), w, r) // Passing context to getUsername
	})
	http.HandleFunc("/changeUsername", func(w http.ResponseWriter, r *http.Request) {
		server.changeUsernameHandler(w, r) // Handling the change username request
	})

	http.HandleFunc("/clearHistory", func(w http.ResponseWriter, r *http.Request) {
		server.clearChatHistoryHTTPHandler(r.Context(), w, r)
	})

	http.HandleFunc("/clearAll", func(w http.ResponseWriter, r *http.Request) {
		server.clearAllHTTPHandler(r.Context(), w, r)
	})

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	log.Println("Starting server on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", corsHandler(http.DefaultServeMux)))
}
