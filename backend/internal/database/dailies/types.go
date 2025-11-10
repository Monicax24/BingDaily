package dailies

import "time"

type Daily struct {
	PostID      string    `json:"post_id"`
	CommunityID string    `json:"community_id"`
	Picture     string    `json:"picture"`
	Caption     string    `json:"caption"`
	Author      string    `json:"author"`
	TimePosted  time.Time `json:"time_posted"`
	Likes       []string  `json:"likes"`
}
