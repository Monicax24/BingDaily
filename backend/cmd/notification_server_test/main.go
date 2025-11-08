package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	push "bingdaily/backend/internal/firebase"
)

type registerReq struct {
	UserID   string `json:"userId"`
	Token    string `json:"token"`
	Platform string `json:"platform,omitempty"`
}

type notifyUserReq struct {
	UserID   string            `json:"userId"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Category string            `json:"category,omitempty"`
	ImageURL string            `json:"imageUrl,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

type notifyTokenReq struct {
	Token    string            `json:"token"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	ImageURL string            `json:"imageUrl,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

type batchNotifyReq struct {
	UserIDs  []string          `json:"userIds"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Category string            `json:"category,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

type topicNotifyReq struct {
	Topic    string            `json:"topic"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Category string            `json:"category,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

type topicSubscribeReq struct {
	Tokens []string `json:"tokens"`
	Topic  string   `json:"topic"`
}

type preferencesReq struct {
	Preferences push.NotificationPreferences `json:"preferences"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func registerToken(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Token == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	push.SaveToken(req.UserID, req.Token, req.Platform)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func unregisterToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	push.RemoveToken(req.Token)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func notifyUser(w http.ResponseWriter, r *http.Request) {
	var req notifyUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Title == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	successCount, invalidTokens, err := push.SendToUserWithPreferences(
		context.Background(),
		req.UserID,
		req.Title,
		req.Body,
		req.Category,
		req.Data,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"sent":          successCount,
		"invalidTokens": len(invalidTokens),
	})
}

func notifyToken(w http.ResponseWriter, r *http.Request) {
	var req notifyTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" || req.Title == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var id string
	var err error

	if req.ImageURL != "" {
		id, err = push.SendNotificationWithImage(context.Background(), req.Title, req.Body, req.Token, req.ImageURL, req.Data)
	} else {
		id, err = push.SendNotification(context.Background(), req.Title, req.Body, req.Token, req.Data)
	}

	if err != nil {
		http.Error(w, "send failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func notifyBatch(w http.ResponseWriter, r *http.Request) {
	var req batchNotifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.UserIDs) == 0 || req.Title == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	totalSent := 0
	totalFailed := 0

	for _, userID := range req.UserIDs {
		successCount, _, err := push.SendToUserWithPreferences(
			context.Background(),
			userID,
			req.Title,
			req.Body,
			req.Category,
			req.Data,
		)
		if err != nil {
			totalFailed++
			continue
		}
		totalSent += successCount
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"totalUsers": len(req.UserIDs),
		"sent":       totalSent,
		"failed":     totalFailed,
	})
}

func notifyTopic(w http.ResponseWriter, r *http.Request) {
	var req topicNotifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Topic == "" || req.Title == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Data == nil {
		req.Data = make(map[string]string)
	}
	req.Data["category"] = req.Category

	id, err := push.SendToTopic(context.Background(), req.Topic, req.Title, req.Body, req.Data)
	if err != nil {
		http.Error(w, "send failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func subscribeToTopic(w http.ResponseWriter, r *http.Request) {
	var req topicSubscribeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Topic == "" || len(req.Tokens) == 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	resp, err := push.SubscribeToTopic(context.Background(), req.Tokens, req.Topic)
	if err != nil {
		http.Error(w, "subscribe failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{
		"successCount": resp.SuccessCount,
		"failureCount": resp.FailureCount,
	})
}

func unsubscribeFromTopic(w http.ResponseWriter, r *http.Request) {
	var req topicSubscribeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Topic == "" || len(req.Tokens) == 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	resp, err := push.UnsubscribeFromTopic(context.Background(), req.Tokens, req.Topic)
	if err != nil {
		http.Error(w, "unsubscribe failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{
		"successCount": resp.SuccessCount,
		"failureCount": resp.FailureCount,
	})
}

func updatePreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	var req preferencesReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	push.SaveUserPreferences(userID, req.Preferences)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func getPreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	prefs := push.GetUserPreferences(userID)
	writeJSON(w, http.StatusOK, prefs)
}

func getUserTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	records := push.GetTokenRecords(userID)
	writeJSON(w, http.StatusOK, map[string]any{"tokens": records})
}

func getNotificationHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs := push.GetNotificationLogs(userID, limit)
	writeJSON(w, http.StatusOK, map[string]any{"logs": logs})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func main() {
	// Token management
	http.HandleFunc("/register-token", registerToken)
	http.HandleFunc("/unregister-token", unregisterToken)
	http.HandleFunc("/user/tokens", getUserTokens)

	// Sending notifications
	http.HandleFunc("/notify/user", notifyUser)
	http.HandleFunc("/notify/token", notifyToken)
	http.HandleFunc("/notify/batch", notifyBatch)
	http.HandleFunc("/notify/topic", notifyTopic)

	// Topic management
	http.HandleFunc("/topic/subscribe", subscribeToTopic)
	http.HandleFunc("/topic/unsubscribe", unsubscribeFromTopic)

	// Preferences
	http.HandleFunc("/preferences", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getPreferences(w, r)
		} else if r.Method == http.MethodPut {
			updatePreferences(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// History
	http.HandleFunc("/notifications/history", getNotificationHistory)

	// Health check
	http.HandleFunc("/health", healthCheck)

	// Background task: cleanup stale tokens every 24 hours
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			removed := push.CleanupStaleTokens(90 * 24 * time.Hour) // 90 days
			log.Printf("Cleaned up %d stale tokens", removed)
		}
	}()

	log.Println("🚀 Notifications server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
