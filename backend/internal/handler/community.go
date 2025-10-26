package handler

import (
	"github.com/gin-gonic/gin"
)

func CommunityRoutes(router *gin.Engine) {
	communityGroup := router.Group("/community")
	communityGroup.GET("/:communityId", fetchCommunityData)
}

func fetchCommunityData(c *gin.Context) {
	c.JSON(200, gin.H{"status": "comm data fetched"})
}
