package communities

type Community struct {
	CommunityID   string   `json:"community_id"`
	Picture       string   `json:"picture"`
	Description   string   `json:"description"`
	Members       []string `json:"members"`   // Array of user IDs
	Posts         []string `json:"posts"`     // Array of post IDs
	PostTime      string   `json:"post_time"` // Default time for daily posts (e.g., "09:00")
	DefaultPrompt string   `json:"default_prompt"`
	Name          string   `json:"name"`
}
