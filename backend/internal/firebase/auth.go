package firebase

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func InitializeAuthClient() *auth.Client {
	// check to see if env variables set
	fmt.Println(os.Getenv("PG_DSN"))

	if keyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); keyPath == "" {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS env variable not set!")
		os.Exit(1)
	}

	// authenticate to AdminSDK
	config := &firebase.Config{
		ProjectID: "bing-daily",
	}
	app, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		fmt.Println("Error initializing app")
	}

	// create authentication service
	auth, err := app.Auth(context.TODO())
	if err != nil {
		fmt.Println("Error initializing auth client")
		os.Exit(1)
	}

	return auth
}

func DecodeToken(token string, auth *auth.Client) string {
	uid, err := auth.VerifyIDTokenAndCheckRevoked(context.TODO(), token)
	if err != nil {
		return ""
	}
	return uid.UID
}
