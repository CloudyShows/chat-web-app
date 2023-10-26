// redis.go
package main

import (
	"log"
)

func clearChatHistory() {
	err := rdb.Del(ctx, "chatHistory").Err()
	if err != nil {
		log.Println("Error clearing chat history:", err)
	} else {
		log.Println("Chat history cleared.")
	}
}

func updateUsernameInRedis(clientIP string, username string) error {
	err := rdb.Set(ctx, "username:"+clientIP, username, 0).Err()
	if err != nil {
		log.Println("Error saving username:", err)
		return err
	}
	// New Debugging code
	retrievedUsername, err := rdb.Get(ctx, "username:"+clientIP).Result()
	if err != nil {
		log.Println("Error retrieving after saving:", err)
	}
	log.Println("Just saved and retrieved:", retrievedUsername)
	return nil
}
