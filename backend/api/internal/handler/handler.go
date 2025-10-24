package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine){
	router.GET("/user/get", getUser)
}


func getUser(c *gin.Context) {
	c.JSON(200, gin.H{"user": "user1",})
}
