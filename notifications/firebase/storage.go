package firebase

import (
	"sync"
	"time"
)

type TokenRecord struct {
	Token     string
	Platform  string // "ios" | "android" | "web"
	UpdatedAt time.Time
}

type NotificationPreferences struct {
	Enabled  bool
	Channels struct {
		Push  bool
		Email bool
		SMS   bool
	}
	Categories map[string]bool // marketing, transactional, alerts, etc.
}

func (p NotificationPreferences) IsCategoryEnabled(category string) bool {
	if len(p.Categories) == 0 {
		return true // default to enabled if no categories set
	}
	enabled, exists := p.Categories[category]
	return exists && enabled
}

type NotificationLog struct {
	UserID       string
	Title        string
	Body         string
	Category     string
	SuccessCount int
	FailureCount int
	Timestamp    time.Time
}

var (
	mu               sync.RWMutex
	userTokens       = map[string]map[string]TokenRecord{}  // userID -> token -> record
	userPreferences  = map[string]NotificationPreferences{} // userID -> preferences
	notificationLogs = []NotificationLog{}                  // in-memory logs (use DB in production)
)

// SaveToken upserts a token for a user
func SaveToken(userID, token, platform string) {
	mu.Lock()
	defer mu.Unlock()
	if userTokens[userID] == nil {
		userTokens[userID] = map[string]TokenRecord{}
	}
	userTokens[userID][token] = TokenRecord{
		Token:     token,
		Platform:  platform,
		UpdatedAt: time.Now(),
	}
}

// RemoveToken deletes a token (e.g., logout/uninstall)
func RemoveToken(token string) {
	mu.Lock()
	defer mu.Unlock()
	for uid := range userTokens {
		delete(userTokens[uid], token)
		if len(userTokens[uid]) == 0 {
			delete(userTokens, uid)
		}
	}
}

// RemoveUserToken removes a specific token for a user
func RemoveUserToken(userID, token string) {
	mu.Lock()
	defer mu.Unlock()
	if userTokens[userID] != nil {
		delete(userTokens[userID], token)
		if len(userTokens[userID]) == 0 {
			delete(userTokens, userID)
		}
	}
}

// GetTokens returns all tokens for a user
func GetTokens(userID string) []string {
	mu.RLock()
	defer mu.RUnlock()
	var out []string
	for t := range userTokens[userID] {
		out = append(out, t)
	}
	return out
}

// GetTokenRecords returns all token records for a user
func GetTokenRecords(userID string) []TokenRecord {
	mu.RLock()
	defer mu.RUnlock()
	var out []TokenRecord
	for _, record := range userTokens[userID] {
		out = append(out, record)
	}
	return out
}

// SaveUserPreferences saves notification preferences for a user
func SaveUserPreferences(userID string, prefs NotificationPreferences) {
	mu.Lock()
	defer mu.Unlock()
	userPreferences[userID] = prefs
}

// GetUserPreferences returns notification preferences for a user
func GetUserPreferences(userID string) NotificationPreferences {
	mu.RLock()
	defer mu.RUnlock()
	prefs, exists := userPreferences[userID]
	if !exists {
		// Return default preferences
		return NotificationPreferences{
			Enabled: true,
			Channels: struct {
				Push  bool
				Email bool
				SMS   bool
			}{Push: true, Email: true, SMS: false},
			Categories: make(map[string]bool),
		}
	}
	return prefs
}

// LogNotification logs a notification send attempt
func LogNotification(userID, title, body, category string, successCount, failureCount int) {
	mu.Lock()
	defer mu.Unlock()
	notificationLogs = append(notificationLogs, NotificationLog{
		UserID:       userID,
		Title:        title,
		Body:         body,
		Category:     category,
		SuccessCount: successCount,
		FailureCount: failureCount,
		Timestamp:    time.Now(),
	})
}

// GetNotificationLogs returns logs for a user
func GetNotificationLogs(userID string, limit int) []NotificationLog {
	mu.RLock()
	defer mu.RUnlock()
	var logs []NotificationLog
	count := 0
	// Iterate in reverse to get most recent first
	for i := len(notificationLogs) - 1; i >= 0 && count < limit; i-- {
		if notificationLogs[i].UserID == userID {
			logs = append(logs, notificationLogs[i])
			count++
		}
	}
	return logs
}

// CleanupStaleTokens removes tokens older than the specified duration
func CleanupStaleTokens(maxAge time.Duration) int {
	mu.Lock()
	defer mu.Unlock()

	removed := 0
	cutoff := time.Now().Add(-maxAge)

	for userID, tokens := range userTokens {
		for token, record := range tokens {
			if record.UpdatedAt.Before(cutoff) {
				delete(tokens, token)
				removed++
			}
		}
		if len(tokens) == 0 {
			delete(userTokens, userID)
		}
	}

	return removed
}
