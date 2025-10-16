package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	push "myfirebaseapp/notifications/firebase"
)

type NotificationRequest struct {
	Token string            `json:"token"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	id, err := push.SendNotification(context.Background(), req.Title, req.Body, req.Token, req.Data)
	if err != nil {
		http.Error(w, "Failed to send notification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Sent notification with ID: " + id))
}

func main() {
	http.HandleFunc("/notify", notifyHandler)
	log.Println("🚀 Notifications server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
