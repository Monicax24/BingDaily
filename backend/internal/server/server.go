package server

import (
	"bingdaily/backend/internal/firebase"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"firebase.google.com/go/v4/auth"
)

type Server struct {
	Router     *gin.Engine
	DB         *sql.DB
	AuthClient *auth.Client
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

// TODO: trace the middleware and figure out actual execution flow
// Register all the different routes
func RegisterRoutes(s *Server) {
	s.Router.Use(s.authenticateUser) // for now require auth by default
	s.Router.Use(s.errorHandling)

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

// Authentication middleware
func (s *Server) authenticateUser(c *gin.Context) {
	// check for auth header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Error(errors.New("authorization header missing"))
		return
	}
	// check if valid format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.Error(errors.New("invalid authentication header format"))
		return
	}
	// decode token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	uid := firebase.DecodeToken(token, s.AuthClient)
	if uid == "" {
		c.Error(errors.New("invalid token"))
		return
	}
	c.Set("userId", uid)
}

// TODO: come up with better error handling scheme (panics, dont leak internal, etc.)
// Error handling middleware
func (s *Server) errorHandling(c *gin.Context) {
	// BUG: not guarenteed that auth failure will call this
	// if there was an error during auth (if there was auth)
	if c.Errors.Last() != nil {
		err := c.Errors.Last().Error()
		sendReponse(c, false, err, nil)
		c.Abort()
		return
	}

	// process rest of request
	c.Next()

	// if there was an error during processing
	if c.Errors.Last() != nil {
		err := c.Errors.Last().Error()
		sendReponse(c, false, err, nil)
		c.Abort()
		return
	}
}
