package main

//change to send to everyone on the app
//send text messages based on subscriptions
//default -> 15 min to make sure they keep their streak
//community notifications for new posts, replies, likes, follows, etc. -> depending on what they subscribe to

//need to send messages depnding ont he topic/community users subscribed to
//community_postReminder
//community_endPromptReminder

import (
	"context"
	"fmt"
	"log"

	push "bingdaily/backend/internal/firebase"
)

func main() {
	// Replace with a real FCM token when ready
	//token := "fyCrqDAdTruYXLIC1shpEd:APA91bFdjTWf1xQVi2XVcGf8g9senyILS2pk1TH1joyr2rLqzFMpV-6I592H5puJOzJAa6Y6bxXmVmJYJHSe9dJVMzL_91OCdRabtDjem5TPtnDcmXRgLeQ"
	token := "cTagTvrjU0uHvrd0P7ZDDD:APA91bEfLDa54VI-EjA4Ds60ScHa2sxUxD_pV_KL9JyvW5qhUJ6yv7Ml6eJJfJORuxSjBqza437DPj148pjDbzh3HWF_z9ToFTIdTO_qRbPvI1VhjUi3MGc"

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
