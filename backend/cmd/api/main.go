package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"bingdaily/backend/internal/database"
	"bingdaily/backend/internal/firebase"
	"bingdaily/backend/internal/server"
	"bingdaily/backend/internal/storage"
)

func main() {
	// ensure that env var is set
	env_set := os.Getenv("BINGDAILY_ENV_SET")

	if env_set == "" {
		fmt.Println("Environment variables are not set!")
		os.Exit(1)
	}

	// init storage
	strge := storage.InitializeStorage()

	// init auth
	authClient := firebase.InitializeAuthClient()

	// init database
	db := database.InitializeDatabase()

	// init router
	router := gin.Default()

	// create and run the server
	s := &server.Server{
		AuthClient: authClient,
		DB:         db,
		Router:     router,
		Storage:    strge,
	}

	server.RegisterRoutes(s)

	s.Router.Run()
}
