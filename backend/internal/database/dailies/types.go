package dailies

import "time"

type Daily struct {
	PostID      int       `json:"post_id"`
	CommunityID int       `json:"community_id"`
	Picture     string    `json:"picture"`
	Caption     string    `json:"caption"`
	Author      int       `json:"author"`
	TimePosted  time.Time `json:"time_posted"`
	Likes       int       `json:"likes"` // Changed from []int to int (count only)
}
