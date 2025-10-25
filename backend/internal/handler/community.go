package handler

import (
	"github.com/gin-gonic/gin"
)

type Community struct {
	CommunityId string `json:"communityId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`

	MemberCnt int `json:"memberCnt"`
}

func CommunityRoutes(router *gin.Engine) {
	communityGroup := router.Group("/community")
	communityGroup.GET("/:communityId", fetchCommunityData)
}

func fetchCommunityData(c *gin.Context) {
	c.JSON(200, gin.H{"status": "comm data fetched"})
}
