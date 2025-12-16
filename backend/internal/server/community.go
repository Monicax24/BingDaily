package server

import (
	"bingdaily/backend/internal/database/communities"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Community struct {
	CommunityId string `json:"communityId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`

	MemberCnt int `json:"memberCnt"`
}

func (s *Server) fetchCommunityData(c *gin.Context) {
	communityId := c.Param("communityId")

	// check to see if community exists
	comm, err := communities.GetCommunity(
		s.DB,
		communityId,
	)
	if err != nil {
		fmt.Printf("Error Code: %v\n", err)
		sendResponse(c, false, "invalid community id", nil)
		return
	}

	// if it does exist, send back to user
	commResponse := &Community{
		CommunityId: comm.CommunityID,
		Name:        comm.Name,
		Description: comm.Description,
		Prompt:      comm.Prompt,
		MemberCnt:   len(comm.Members),
	}

	sendResponse(c, true, fmt.Sprintf("retreived %s", communityId), commResponse)
}
