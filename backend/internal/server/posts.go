package server

import (
	"bingdaily/backend/internal/database/communities"
	"bingdaily/backend/internal/database/dailies"
	"bingdaily/backend/internal/storage"
	"fmt"

	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	PostId      string    `json:"postId"`
	CommunityId string    `json:"communityId"`
	UserId      string    `json:"author"`
	Caption     string    `json:"caption"`
	ImageUrl    string    `json:"imageUrl"`
	TimePosted  time.Time `json:"timePosted"`
}

type FetchCommunityPostsResponse struct {
	Posts []*Post `json:"posts"`
}

type CreatePostRequest struct {
	CommunityId string `json:"communityId" binding:"required"`
	Caption     string `json:"caption" binding:"required"`
}

type CreatePostResponse struct {
	PostId    string `json:"postId" binding:"required"`
	UploadUrl string `json:"uploadUrl" binding:"required"`
}

func (s *Server) fetchCommunityPosts(c *gin.Context) {
	communityId := c.Param("communityId")
	userId := c.Value("userId").(string)

	// check to see if community exists
	comm, err := communities.GetCommunity(
		s.DB,
		communityId,
	)
	if err != nil {
		sendResponse(c, false, "invalid community id", nil)
		return
	}

	// check to see if user in community
	if !slices.Contains(comm.Members, userId) {
		sendResponse(c, false, "unauthorized operation", nil)
		return
	}

	// retrieve all the posts
	dlies, err := dailies.FetchDailiesFromCommunity(s.DB, communityId)
	if err != nil {
		sendResponse(c, false, "internal error", nil)
		return
	}
	var posts []*Post
	for _, daily := range dlies {
		// TODO: how can we handle placeholder?
		// add image URL if it exists
		imageUrl := ""
		if daily.Picture != "" {
			url, err := s.Storage.GenerateDownloadURL(storage.POST_PICTURES, daily.Picture)
			// error getting picture so send placeholder
			// TODO: send local placeholder, not publicy-hosted
			if err != nil {
				imageUrl = "https://www.setra.com/hubfs/Sajni/crc_error.jpg"
			} else {
				imageUrl = url
			}
		}
		res := &Post{
			PostId:      daily.PostID,
			CommunityId: daily.CommunityID,
			UserId:      daily.Author,
			Caption:     daily.Caption,
			TimePosted:  daily.TimePosted,
			ImageUrl:    imageUrl,
		}
		posts = append(posts, res)
	}
	res := &FetchCommunityPostsResponse{
		Posts: posts,
	}
	sendResponse(c, true, fmt.Sprintf("fetched %d posts", len(posts)), res)
}

// TODO: right now there are 3 calls to DB, minimize # of calls
func (s *Server) uploadPost(c *gin.Context) {
	req := &CreatePostRequest{}
	userId := c.Value("userId").(string)

	// check if valid request
	err := c.ShouldBind(&req)
	if err != nil {
		sendResponse(c, false, "invalid request body", nil)
		return
	}

	// check if user in community
	in, err := communities.UserInCommunity(s.DB, req.CommunityId, userId)
	if !in {
		sendResponse(c, false, "unauthorized operation", nil)
		return
	} else if err != nil {
		sendResponse(c, false, "internal error", nil)
		return
	}

	// check if user has posted before
	posted, err := dailies.HasPostedToday(s.DB, userId, req.CommunityId)
	if err != nil {
		sendResponse(c, false, "internal error", nil)
		return
	} else if posted {
		sendResponse(c, false, "already posted", nil)
		return
	}

	// TODO: s3 + sql operations can be done async?
	// create picture
	pictureId := storage.CreatePictureId()
	uploadUrl, err := s.Storage.GenerateUploadURL(storage.POST_PICTURES, pictureId)
	// add to db
	dailyId, err1 := dailies.CreateDaily(
		s.DB,
		req.CommunityId,
		pictureId,
		req.Caption,
		userId,
	)
	if err != nil || err1 != nil {
		sendResponse(c, false, "internal error", nil)
		return
	}

	res := &CreatePostResponse{
		PostId:    dailyId,
		UploadUrl: uploadUrl,
	}
	sendResponse(c, true, "post uploaded", res)
}
