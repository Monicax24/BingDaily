package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) fetchCommunityPosts(c *gin.Context) {
	communityId := c.Param("communityId")
	c.JSON(200, gin.H{"status": communityId})
}

func (s *Server) uploadPost(c *gin.Context) {
	c.JSON(200, gin.H{"status": "post uploaded"})
}
