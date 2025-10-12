package main

import (
	"context"
	"fmt"
	"log"

	push "myfirebaseapp/backend/firebase"
)

func main() {
	// Replace with a real FCM token when ready
	token := "DEVICE_TOKEN"

	id, err := push.SendNotification(context.Background(),
		"Test Title",
		"Test Body",
		token,
		map[string]string{"foo": "bar"}, // or nil
	)
	if err != nil {
		log.Fatalf("failed to send: %v", err)
	}
	fmt.Println("✅ Sent message:", id)
}
