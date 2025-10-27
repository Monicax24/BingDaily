package firebase

import (
	"sync"
	"time"
)

type Metrics struct {
	mu                sync.RWMutex
	NotificationsSent int64
	NotificationsFail int64
	TokensRegistered  int64
	TokensRemoved     int64
	LastReset         time.Time
}

var globalMetrics = &Metrics{
	LastReset: time.Now(),
}

func IncrementNotificationsSent(count int) {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	globalMetrics.NotificationsSent += int64(count)
}

func IncrementNotificationsFailed(count int) {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	globalMetrics.NotificationsFail += int64(count)
}

func IncrementTokensRegistered() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	globalMetrics.TokensRegistered++
}

func IncrementTokensRemoved() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	globalMetrics.TokensRemoved++
}

func GetMetrics() Metrics {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()
	return *globalMetrics
}

func ResetMetrics() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	globalMetrics.NotificationsSent = 0
	globalMetrics.NotificationsFail = 0
	globalMetrics.TokensRegistered = 0
	globalMetrics.TokensRemoved = 0
	globalMetrics.LastReset = time.Now()
}
