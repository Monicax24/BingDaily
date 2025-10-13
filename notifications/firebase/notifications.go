package firebasepush

import (
	"context"
	"os"
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
		app, appErr = firebase.NewApp(context.Background(), nil, option.WithCredentialsFile("backend/firebase/serviceAccount.json"))
	})
	return app, appErr
}

// SendNotification sends a notification to a single device token.
// will pass nil for data if you don't need custom key/values.
// still need fcm token
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
	}
	return client.Send(ctx, msg)
}
