package server

import (
	"bingdaily/backend/internal/firebase"
	"bingdaily/backend/internal/storage"
	"errors"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Router     *gin.Engine
	DB         *pgxpool.Pool
	AuthClient *auth.Client
	Storage    *storage.Storage
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
	communityGroup := s.Router.Group("/communities")
	communityGroup.GET("/:communityId", s.fetchCommunityData)

	// post routes
	postsGroup := s.Router.Group("/communities/posts")
	postsGroup.GET("/:communityId", s.fetchCommunityPosts)
	postsGroup.POST("/upload", s.uploadPost)

	// user routes
	userGroup := s.Router.Group("/users")
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

	// WARNING: remove this for production
	// temp access token
	if token == "X7f6sH9nAEU+9U6vzLNGK0EqzgFwALcOdNbpsHwplx3E04488E12QA=" {
		c.Set("userId", "697b8a69-0c01-4ccb-aabc-6dccd6a22fa3")
		return
	}

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
