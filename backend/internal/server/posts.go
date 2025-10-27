package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	PostId      string `json:"postId"`
	CommunityId string `json:"communityId"`
	UserId      string `json:"author"`

	Caption    string    `json:"caption"`
	TimePosted time.Time `json:"timePosted"`
}

var fakePost1 = &Post{
	PostId:      "fakepost1id",
	CommunityId: "fakecommunityid",
	UserId:      "fakeuserid",

	Caption:    "This is a test caption for test post 1!",
	TimePosted: time.Now(),
}

var fakePost2 = &Post{
	PostId:      "fakepost2id",
	CommunityId: "fakecommunityid",
	UserId:      "fakeuserid",

	Caption:    "This is a test caption for test post 2!",
	TimePosted: time.Now(),
}

func (s *Server) fetchCommunityPosts(c *gin.Context) {
	communityId := c.Param("communityId")

	arr := [2]*Post{fakePost1, fakePost2}
	sendReponse(
		c,
		true,
		fmt.Sprintf("Retrieved posts for %s successfully!", communityId),
		arr,
	)
}

func (s *Server) uploadPost(c *gin.Context) {
	sendReponse(c, true, "Post uploaded successfully", nil)
}
