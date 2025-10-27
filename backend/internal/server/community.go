package server

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

var fakeComm *Community = &Community{
	CommunityId: "fakecommunityid",
	Name:        "Test Community",
	Description: "This community is a test community",
	Prompt:      "This is a test prompt",
	MemberCnt:   67,
}

func (s *Server) fetchCommunityData(c *gin.Context) {
	sendReponse(c, true, "", fakeComm)
}
