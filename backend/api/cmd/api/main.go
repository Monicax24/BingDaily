package main

import (
	"bingdaily/api/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	handler.RegisterRoutes(router)
	router.Run()
}
