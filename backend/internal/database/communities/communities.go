package communities

import (
	"database/sql"

	"github.com/google/uuid"
)

type Community struct {
	CommunityID string   `json:"community_id"`
	Picture     string   `json:"picture"`
	Description string   `json:"description"`
	Members     []string `json:"members"`   // Array of user IDs
	Posts       []string `json:"posts"`     // Array of post IDs
	PostTime    string   `json:"post_time"` // Default time for daily posts (e.g., "09:00")
	Prompt      string   `json:"prompt"`
	Name        string   `json:"name"`
}

// Create a new community
func CreateCommunity(db *sql.DB, name, picture, description, postTime, prompt string) (string, error) {
	// Generate a unique community ID
	communityID := uuid.New().String()

	_, err := db.Exec(`
		INSERT INTO communities (community_id, name, picture, description, members, posts, post_time, prompt) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		communityID, name, picture, description, "{}", "{}", postTime, prompt)
	return communityID, err
}

// Get community by ID
func GetCommunity(db *sql.DB, communityID string) (*Community, error) {
	var community Community

	err := db.QueryRow(`
		SELECT community_id, picture, description, members, posts, post_time, prompt, name
		FROM communities 
		WHERE community_id = $1`,
		communityID).Scan(&community.CommunityID, &community.Picture, &community.Description,
		&community.Members, &community.Posts, &community.PostTime, &community.Prompt, &community.Name)

	if err != nil {
		return nil, err
	}

	return &community, nil
}

// Join a community
func JoinCommunity(db *sql.DB, userID, communityID string) error {
	_, err := db.Exec(`
		UPDATE communities 
		SET members = array_append(members, $1)
		WHERE community_id = $2 AND NOT $1 = ANY(members)`,
		userID, communityID)

	if err == nil {
		// Also add community to user's communities array
		_, err = db.Exec(`
			UPDATE users 
			SET communities = array_append(communities, $1)
			WHERE user_id = $2 AND NOT $1 = ANY(communities)`,
			communityID, userID)
	}

	return err
}

// Leave a community
func LeaveCommunity(db *sql.DB, userID, communityID string) error {
	_, err := db.Exec(`
		UPDATE communities 
		SET members = array_remove(members, $1)
		WHERE community_id = $2`,
		userID, communityID)

	if err == nil {
		// Also remove community from user's communities array
		_, err = db.Exec(`
			UPDATE users 
			SET communities = array_remove(communities, $1)
			WHERE user_id = $2`,
			communityID, userID)
	}

	return err
}
