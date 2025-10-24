package handler

import (
	"github.com/gin-gonic/gin"
)

// Register all the different routes
func RegisterRoutes(router *gin.Engine) {
	router.GET("/user/register", registerUser)
}
