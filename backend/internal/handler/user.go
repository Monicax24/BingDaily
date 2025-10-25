package handler

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/user")
	userGroup.POST("/register", registerUser)
	userGroup.POST("/update", updateUserProfile)
	userGroup.GET("/:userId", fetchUserProfile)
}

func registerUser(c *gin.Context) {
	c.JSON(200, gin.H{"status": "registered"})
}

func updateUserProfile(c *gin.Context) {
	c.JSON(200, gin.H{"status": "updated"})
}

func fetchUserProfile(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(200, gin.H{"status": userId})
}
