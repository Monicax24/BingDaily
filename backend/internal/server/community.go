package server

import (
	"bingdaily/backend/internal/database/communities"

	"github.com/gin-gonic/gin"

	"fmt"
	"strconv" // temporary until DB fixed
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
	communityId := c.Param("communityId")

	// temp code
	converted, err := strconv.Atoi(communityId)
	if err != nil {
		fmt.Println("Invalid conversion")
	}

	communities.GetCommunity(s.DB, converted)

	sendReponse(c, true, "", fakeComm)
}
