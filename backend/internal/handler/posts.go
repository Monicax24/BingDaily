package handler

import (
	"github.com/gin-gonic/gin"
)

type Post struct {
}

func PostRoutes(router *gin.Engine) {
	postsGroup := router.Group("/community/posts")
	postsGroup.GET("/:communityId", fetchCommunityPosts)
	postsGroup.POST("/upload", uploadPost)
}

func fetchCommunityPosts(c *gin.Context) {
	communityId := c.Param("communityId")
	c.JSON(200, gin.H{"status": communityId})
}

func uploadPost(c *gin.Context) {
	c.JSON(200, gin.H{"status": "post uploaded"})
}
