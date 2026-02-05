package server

import (
	"bingdaily/backend/internal/database/dailies"
	"bingdaily/backend/internal/database/users"
	"bingdaily/backend/internal/firebase"
	"bingdaily/backend/internal/storage"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserId         string    `json:"userId"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	JoinDate       time.Time `json:"joinDate"`
	ProfilePicture string    `json:"profilePicture"`
	Communities    []string  `json:"communities"`
}

type RegisterUserRequest struct {
	Email         string `json:"email" binding:"required"`
	Username      string `json:"username" binding:"required" `
	UpdatePicture bool   `json:"updatePicture"` // no "required" needed since bool
}

type RegisterUserResponse struct {
	UploadUrl string `json:"uploadUrl"`
}

type FetchUserProfileResponse struct {
	User *User `json:"user"`
}

func (s *Server) registerUser(c *gin.Context) {
	// retrieve authenticated user
	userId := c.Value("userId").(string)

	// check if valid request body
	req := &RegisterUserRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		sendResponse(c, false, "invalid request body", nil)
		return
	}

	// check if valid email
	validEmail, err := firebase.ValidEmail(s.AuthClient, userId, req.Email)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		sendResponse(c, false, "internal error", nil)
		return
	} else if !validEmail {
		sendResponse(c, false, "invalid email address", nil)
		return
	}

	// TODO: manage placeholder
	var photoId string = ""
	var data interface{} = nil
	if req.UpdatePicture {
		photoId = storage.CreatePictureId()
		uploadUrl, err := s.Storage.GenerateUploadURL(storage.PROFILE_PHOTOS, photoId)
		data = &RegisterUserResponse{
			UploadUrl: uploadUrl,
		}
		if err != nil {
			sendResponse(c, false, "internal error", nil)
			return
		}
	}

	// register the user (checks for pre-existence)
	registerSuccess, err := users.Register(s.DB, userId, req.Username, req.Email, photoId)
	if err != nil || !registerSuccess {
		sendResponse(c, false, "internal error", nil)
		return
	}
	sendResponse(
		c,
		true,
		"user registered",
		data,
	)
}

func (s *Server) updateUserProfile(c *gin.Context) {
	sendResponse(
		c,
		true,
		"profile updated",
		nil,
	)
}

func (s *Server) fetchUserProfile(c *gin.Context) {
	userId := c.Value("userId").(string)
	dbuser, err := users.GetUser(s.DB, userId)
	if err != nil {
		sendResponse(c, false, "error fetching user", nil)
		return
	}
	user := &User{
		UserId:      dbuser.UserID,
		Email:       dbuser.Email,
		Username:    dbuser.Name,
		JoinDate:    dbuser.JoinedDate,
		Communities: dbuser.Communities,
	}
	// TODO: maybe this can go in helper
	// process photo
	if dbuser.ProfilePicture != "" {
		url, err := s.Storage.GenerateDownloadURL(
			storage.PROFILE_PHOTOS,
			dbuser.ProfilePicture,
		)
		if err != nil {
			sendResponse(c, false, "error fetching user photo", nil)
			return

		}
		user.ProfilePicture = url
	} else {
		user.ProfilePicture = "https://wallpapers.com/images/hd/default-profile-picture-placeholder-kal8zbcust2luxh3.jpg"
	}
	res := &FetchUserProfileResponse{User: user}
	sendResponse(
		c,
		true,
		fmt.Sprintf("retrieved %s", userId),
		res,
	)
}

func (s *Server) fetchUserPosts(c *gin.Context) {
	userId := c.Value("userId").(string)
	dlies, err := dailies.FetchDailiesFromUser(s.DB, userId)
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
