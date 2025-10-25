package main

import (
	"bingdaily/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	handler.RegisterRoutes(router)
	router.Run()
}
