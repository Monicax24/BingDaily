package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`

	JoinDate    time.Time `json:"joinDate"`
	Communities []string  `json:"communities"`
}

var fakeUser *User = &User{
	UserId:   "fakeuserid",
	Email:    "fake@email.com",
	Username: "fakeusername",

	JoinDate:    time.Now(),
	Communities: []string{"fakecommunityid1", "fakecommunityid2"},
}

func (s *Server) registerUser(c *gin.Context) {
	sendReponse(
		c,
		true,
		"User registration success!",
		nil,
	)
}

func (s *Server) updateUserProfile(c *gin.Context) {
	sendReponse(
		c,
		true,
		"User profile update success!",
		nil,
	)
}

func (s *Server) fetchUserProfile(c *gin.Context) {
	userId := c.Param("userId")
	sendReponse(
		c,
		true,
		fmt.Sprintf("Retrieved user %s successfully!", userId),
		fakeUser,
	)
}
