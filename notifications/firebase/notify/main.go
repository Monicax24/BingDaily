package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	push "myfirebaseapp/notifications/firebase" // <-- update to your module path
)

type registerReq struct {
	UserID   string `json:"userId"`
	Token    string `json:"token"`
	Platform string `json:"platform,omitempty"`
}

type notifyUserReq struct {
	UserID string            `json:"userId"`
	Title  string            `json:"title"`
	Body   string            `json:"body"`
	Data   map[string]string `json:"data,omitempty"`
}

type notifyTokenReq struct {
	Token string            `json:"token"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
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
	writeJSON(w, http.StatusNoContent, map[string]string{"status": "ok"})
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
	writeJSON(w, http.StatusNoContent, map[string]string{"status": "ok"})
}

func notifyUser(w http.ResponseWriter, r *http.Request) {
	var req notifyUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Title == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	tokens := push.GetTokens(req.UserID)
	if len(tokens) == 0 {
		http.Error(w, "no tokens for user", http.StatusNotFound)
		return
	}
	// naive fan-out for MVP
	var ids []string
	for _, t := range tokens {
		id, err := push.SendNotification(context.Background(), req.Title, req.Body, t, req.Data)
		if err != nil {
			// In production: detect "not-registered" and call push.RemoveToken(t)
			http.Error(w, "send failed: "+err.Error(), http.StatusBadGateway)
			return
		}
		ids = append(ids, id)
	}
	writeJSON(w, http.StatusOK, map[string]any{"sent": len(ids), "ids": ids})
}

func notifyToken(w http.ResponseWriter, r *http.Request) {
	var req notifyTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" || req.Title == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id, err := push.SendNotification(context.Background(), req.Title, req.Body, req.Token, req.Data)
	if err != nil {
		http.Error(w, "send failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func main() {
	http.HandleFunc("/register-token", registerToken)
	http.HandleFunc("/unregister-token", unregisterToken)
	http.HandleFunc("/notify/user", notifyUser)
	http.HandleFunc("/notify/token", notifyToken)

	log.Println("🚀 Notifications server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
