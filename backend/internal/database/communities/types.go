package communities

type Community struct {
	CommunityID   int    `json:"community_id"` // Changed from string to int
	Picture       string `json:"picture"`
	Description   string `json:"description"`
	Members       []int  `json:"members"`    // Array of user IDs
	Moderators    []int  `json:"moderators"` // Array of user IDs
	Posts         []int  `json:"posts"`      // Array of post IDs
	PostTime      string `json:"post_time"`  // Default time for daily posts (e.g., "09:00")
	DefaultPrompt string `json:"default_prompt"`
}
