package server

import (
	"bingdaily/backend/internal/database/communities"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Community struct {
	CommunityId string `json:"communityId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
	MemberCnt   int    `json:"memberCnt"`
}

type FetchCommunityDataResponse struct {
	Community *Community `json:"community"`
}

type FetchCommunitiesResponse struct {
	Communities []*Community `json:"communities"`
}

func (s *Server) fetchCommunityData(c *gin.Context) {
	communityId := c.Param("communityId")

	// check to see if community exists
	comm, err := communities.GetCommunity(
		s.DB,
		communityId,
	)
	if err != nil {
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
	res := &FetchCommunityDataResponse{
		Community: commResponse,
	}
	sendResponse(c, true, fmt.Sprintf("retrieved %s", communityId), res)
}

func (s *Server) joinCommunity(c *gin.Context) {
	userId := c.Value("userId").(string)
	communityId := c.Param("communityId")

	err := communities.JoinCommunity(s.DB, userId, communityId)
	if err == pgx.ErrNoRows {
		sendResponse(c, false, "already joined", nil)
		return
	} else if err != nil {
		sendResponse(c, false, "internal error", nil)
		return
	}

	sendResponse(c, true, fmt.Sprintf("joined %s", communityId), nil)
}

func (s *Server) leaveCommunity(c *gin.Context) {
	userId := c.Value("userId").(string)
	communityId := c.Param("communityId")

	err := communities.LeaveCommunity(s.DB, userId, communityId)
	if err == pgx.ErrNoRows {
		sendResponse(c, false, "not in community", nil)
		return
	} else if err != nil {
		sendResponse(c, false, "internal error", nil)
		return
	}
	sendResponse(c, true, fmt.Sprintf("left %s", communityId), nil)
}

func (s *Server) fetchCommunities(c *gin.Context) {
	db_comms, err := communities.ListCommunites(s.DB)
	if err != nil {
		sendResponse(c, false, "internal error", nil)
		return
	}

	var comms []*Community
	for _, comm := range db_comms {
		res := &Community{
			CommunityId: comm.CommunityID,
			Name:        comm.Name,
			Description: comm.Description,
			Prompt:      comm.Description,
			MemberCnt:   comm.MemberCnt,
		}
		comms = append(comms, res)
	}

	res := &FetchCommunitiesResponse{
		Communities: comms,
	}
	sendResponse(c, true, fmt.Sprintf("fetched %d communities", len(comms)), res)
}
