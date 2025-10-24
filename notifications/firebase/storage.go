//handles user token database for firebase notifications

package firebase

import (
	"sync"
	"time"
)

type TokenRecord struct {
	Token     string
	Platform  string // "ios" | "android" | "web" (optional)
	UpdatedAt time.Time
}

var (
	mu         sync.RWMutex
	userTokens = map[string]map[string]TokenRecord{} // userID -> token -> record
)

// SaveToken upserts a token for a user.
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

// RemoveToken deletes a token (e.g., logout/uninstall).
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

// GetTokens returns all tokens for a user.
func GetTokens(userID string) []string {
	mu.RLock()
	defer mu.RUnlock()
	var out []string
	for t := range userTokens[userID] {
		out = append(out, t)
	}
	return out
}
