// main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	http.HandleFunc("/ws", handler)
	http.HandleFunc("/getUsername", getUsername)
	http.HandleFunc("/clear", clearChatHistoryHTTPHandler)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	log.Println("Starting server on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", corsHandler(http.DefaultServeMux)))
}
