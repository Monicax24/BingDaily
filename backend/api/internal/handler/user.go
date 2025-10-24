package handler

import "github.com/gin-gonic/gin"

func registerUser(c *gin.Context) {
	c.JSON(200, gin.H{"user": "user1"})
}
