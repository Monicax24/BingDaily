package firebase

// figure who to send notification on

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var (
	app    *firebase.App
	appErr error
	once   sync.Once
)

func initFirebaseApp() (*firebase.App, error) {
	once.Do(func() {
		if json := os.Getenv("FIREBASE_SERVICE_ACCOUNT"); json != "" {
			app, appErr = firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(json)))
			return
		}
		if path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); path != "" {
			app, appErr = firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(path))
			return
		}
		app, appErr = firebase.NewApp(context.Background(), nil, option.WithCredentialsFile("notifications/firebase/serviceAccount.json"))
	})
	return app, appErr
}

// SendNotification sends a notification to a single device token
func SendNotification(ctx context.Context, title, body, token string, data map[string]string) (string, error) {
	app, err := initFirebaseApp()
	if err != nil {
		return "", err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return "", err
	}
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "default",
				Sound:     "default",
				Priority:  messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					Badge: intPtr(1),
				},
			},
		},
	}
	return client.Send(ctx, msg)
}

// SendNotificationWithImage sends a notification with an image
func SendNotificationWithImage(ctx context.Context, title, body, token, imageURL string, data map[string]string) (string, error) {
	app, err := initFirebaseApp()
	if err != nil {
		return "", err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return "", err
	}
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "default",
				Sound:     "default",
				ImageURL:  imageURL,
				Priority:  messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound:            "default",
					Badge:            intPtr(1),
					MutableContent:   true,
					ContentAvailable: true,
				},
			},
			FCMOptions: &messaging.APNSFCMOptions{
				ImageURL: imageURL,
			},
		},
	}
	return client.Send(ctx, msg)
}

// SendMulticast sends the same notification to multiple tokens
func SendMulticast(ctx context.Context, title, body string, tokens []string, data map[string]string) (*messaging.BatchResponse, error) {
	if len(tokens) == 0 {
		return nil, errors.New("no tokens provided")
	}

	app, err := initFirebaseApp()
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	msg := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "default",
				Sound:     "default",
				Priority:  messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					Badge: intPtr(1),
				},
			},
		},
	}

	return client.SendEachForMulticast(ctx, msg)
}

// SendToTopic sends a notification to all devices subscribed to a topic
func SendToTopic(ctx context.Context, topic, title, body string, data map[string]string) (string, error) {
	app, err := initFirebaseApp()
	if err != nil {
		return "", err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return "", err
	}

	msg := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "default",
				Sound:     "default",
				Priority:  messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					Badge: intPtr(1),
				},
			},
		},
	}

	return client.Send(ctx, msg)
}

// SubscribeToTopic subscribes tokens to a topic
func SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	app, err := initFirebaseApp()
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return client.SubscribeToTopic(ctx, tokens, topic)
}

// UnsubscribeFromTopic unsubscribes tokens from a topic
func UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	app, err := initFirebaseApp()
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return client.UnsubscribeFromTopic(ctx, tokens, topic)
}

// SendToUserWithPreferences sends notification to user respecting their preferences
func SendToUserWithPreferences(ctx context.Context, userID, title, body, category string, data map[string]string) (int, []string, error) {
	// Get user preferences
	prefs := GetUserPreferences(userID)
	if !prefs.Enabled || !prefs.Channels.Push {
		return 0, nil, fmt.Errorf("push notifications disabled for user %s", userID)
	}

	// Check category preference
	if category != "" && !prefs.IsCategoryEnabled(category) {
		return 0, nil, fmt.Errorf("category %s disabled for user %s", category, userID)
	}

	// Get tokens
	tokens := GetTokens(userID)
	if len(tokens) == 0 {
		return 0, nil, fmt.Errorf("no tokens for user %s", userID)
	}

	// Add category to data
	if data == nil {
		data = make(map[string]string)
	}
	data["category"] = category

	// Send to all tokens
	response, err := SendMulticast(ctx, title, body, tokens, data)
	if err != nil {
		return 0, nil, err
	}

	// Handle failed tokens (invalid/unregistered)
	var invalidTokens []string
	for i, resp := range response.Responses {
		if !resp.Success {
			errCode := ""
			if resp.Error != nil {
				errCode = strings.ToLower(resp.Error.Error())
			}
			// Check for token-related errors
			if strings.Contains(errCode, "not-registered") ||
				strings.Contains(errCode, "invalid-registration") ||
				strings.Contains(errCode, "invalid-argument") {
				invalidTokens = append(invalidTokens, tokens[i])
			}
		}
	}

	// Remove invalid tokens
	for _, token := range invalidTokens {
		RemoveToken(token)
	}

	// Log notification
	LogNotification(userID, title, body, category, response.SuccessCount, response.FailureCount)

	return response.SuccessCount, invalidTokens, nil
}

func intPtr(i int) *int {
	return &i
}
