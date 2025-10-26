package handler

import (
	"github.com/gin-gonic/gin"
)

// Register all the different routes
func RegisterRoutes(router *gin.Engine) {
	UserRoutes(router)
	CommunityRoutes(router)
	PostRoutes(router)
}

func authenticateUser() {

}

func sendReponse(c *gin.Context, status bool, message string, data interface{}) {
	response := gin.H{
		"status":  "fail",
		"message": nil,
	}
	if status {
		response["status"] = "success"
	}
	if message != "" {
		response["message"] = message
	}
	if data != nil {
		response["data"] = data
	}
	c.JSON(200, response)
}
