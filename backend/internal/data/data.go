package data

import "time"

// User struct represents a user in the system
type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profile_picture"`
	JoinedDate     time.Time `json:"joined_date"`
	Communities    []int     `json:"communities"` // Array of community IDs
	Friends        []int     `json:"friends"`     // Array of user IDs
}

// Post struct represents a daily post
type Post struct {
	PostID      int       `json:"post_id"`
	CommunityID int       `json:"community_id"`
	Picture     string    `json:"picture"`
	Caption     string    `json:"caption"`
	Author      int       `json:"author"` // User ID
	TimePosted  time.Time `json:"time_posted"`
	Likes       []int     `json:"likes"` // Array of user IDs who liked the post
}

// Community struct represents a community
type Community struct {
	CommunityID   int    `json:"community_id"`
	Picture       string `json:"picture"`
	Description   string `json:"description"`
	Members       []int  `json:"members"`    // Array of user IDs
	Moderators    []int  `json:"moderators"` // Array of user IDs
	Posts         []int  `json:"posts"`      // Array of post IDs
	PostTime      string `json:"post_time"`  // Default time for daily posts (e.g., "09:00:00")
	DefaultPrompt string `json:"default_prompt"`
}
