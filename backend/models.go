// models.go
package main

type ChatMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Text     string `json:"text"`
}