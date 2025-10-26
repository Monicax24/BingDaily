package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"firebase.google.com/go/v4/auth"
)

type Server struct {
	Router     *gin.Engine
	DB         *sql.DB
	AuthClient *auth.Client
}

// Register all the different routes
func RegisterRoutes(s *Server) {
	// community routes
	communityGroup := s.Router.Group("/community")
	communityGroup.GET("/:communityId", s.fetchCommunityData)

	// post routes
	postsGroup := s.Router.Group("/community/posts")
	postsGroup.GET("/:communityId", s.fetchCommunityPosts)
	postsGroup.POST("/upload", s.uploadPost)

	// user routes
	userGroup := s.Router.Group("/user")
	userGroup.POST("/register", s.registerUser)
	userGroup.POST("/update", s.updateUserProfile)
	userGroup.GET("/:userId", s.fetchUserProfile)
}

// Helper function for sending server responses
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
	c.JSON(http.StatusOK, response)
}
