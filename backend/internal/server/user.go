package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) registerUser(c *gin.Context) {
	c.JSON(200, gin.H{"status": "registered"})
}

func (s *Server) updateUserProfile(c *gin.Context) {
	c.JSON(200, gin.H{"status": "updated"})
}

func (s *Server) fetchUserProfile(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(200, gin.H{"status": userId})
}
