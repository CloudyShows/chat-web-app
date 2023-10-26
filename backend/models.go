// models.go
package main

import (
	"time"
)

type ChatMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Timestamp time.Time `json:"timestamp"`

}