package main

import (
	"bingdaily/backend/internal/database"
	"bingdaily/backend/internal/firebase"
	"bingdaily/backend/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
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
	}

	server.RegisterRoutes(s)

	s.Router.Run()
}
