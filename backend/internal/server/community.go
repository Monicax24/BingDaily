package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) fetchCommunityData(c *gin.Context) {
	c.JSON(200, gin.H{"status": "comm data fetched"})
}
