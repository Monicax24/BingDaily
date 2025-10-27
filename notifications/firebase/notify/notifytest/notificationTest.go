package notificationTest

import (
	"context"
	"testing"
	"time"
)

func TestTokenManagement(t *testing.T) {
	// Test SaveToken
	SaveToken("user123", "token123", "ios")
	tokens := GetTokens("user123")
	if len(tokens) != 1 {
		t.Errorf("Expected 1 token, got %d", len(tokens))
	}

	// Test RemoveToken
	RemoveToken("token123")
	tokens = GetTokens("user123")
	if len(tokens) != 0 {
		t.Errorf("Expected 0 tokens, got %d", len(tokens))
	}
}

func TestPreferences(t *testing.T) {
	prefs := NotificationPreferences{
		Enabled: true,
		Channels: struct {
			Push  bool
			Email bool
			SMS   bool
		}{Push: true, Email: false, SMS: false},
		Categories: map[string]bool{
			"alerts":    true,
			"marketing": false,
		},
	}

	SaveUserPreferences("user123", prefs)
	savedPrefs := GetUserPreferences("user123")

	if !savedPrefs.Enabled {
		t.Error("Expected preferences to be enabled")
	}

	if !savedPrefs.IsCategoryEnabled("alerts") {
		t.Error("Expected alerts category to be enabled")
	}

	if savedPrefs.IsCategoryEnabled("marketing") {
		t.Error("Expected marketing category to be disabled")
	}
}

func TestCleanupStaleTokens(t *testing.T) {
	// Add some tokens
	SaveToken("user1", "token1", "ios")
	time.Sleep(10 * time.Millisecond)
	SaveToken("user2", "token2", "android")

	// Cleanup tokens older than 5ms
	removed := CleanupStaleTokens(5 * time.Millisecond)
	if removed != 1 {
		t.Errorf("Expected to remove 1 token, removed %d", removed)
	}
}
