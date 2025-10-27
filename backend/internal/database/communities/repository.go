package communities

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Create a new community
func CreateCommunity(db *sql.DB, picture, description, postTime, defaultPrompt string) (int, error) {
	var communityID int
	err := db.QueryRow(`
		INSERT INTO communities (picture, description, members, moderators, posts, post_time, default_prompt) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING community_id`,
		picture, description, "[]", "[]", "[]", postTime, defaultPrompt).Scan(&communityID)
	return communityID, err
}

// Get community by ID
func GetCommunity(db *sql.DB, communityID int) (*Community, error) {
	var community Community
	var membersJSON, moderatorsJSON, postsJSON string

	err := db.QueryRow(`
		SELECT community_id, picture, description, members, moderators, posts, post_time, default_prompt
		FROM communities 
		WHERE community_id = $1`,
		communityID).Scan(&community.CommunityID, &community.Picture, &community.Description,
		&membersJSON, &moderatorsJSON, &postsJSON, &community.PostTime, &community.DefaultPrompt)

	if err != nil {
		return nil, err
	}

	// Parse JSON arrays
	json.Unmarshal([]byte(membersJSON), &community.Members)
	json.Unmarshal([]byte(moderatorsJSON), &community.Moderators)
	json.Unmarshal([]byte(postsJSON), &community.Posts)

	return &community, nil
}

// Join a community
func JoinCommunity(db *sql.DB, userID, communityID int) error {
	_, err := db.Exec(`
		UPDATE communities 
		SET members = members || $1::jsonb
		WHERE community_id = $2 AND NOT members @> $3::jsonb`,
		fmt.Sprintf("[%d]", userID), communityID, fmt.Sprintf("[%d]", userID))
	return err
}

// Leave a community
func LeaveCommunity(db *sql.DB, userID, communityID int) error {
	_, err := db.Exec(`
		UPDATE communities 
		SET members = members - $1
		WHERE community_id = $2`,
		userID, communityID)
	return err
}

// Add moderator to community
func AddModerator(db *sql.DB, communityID, userID int) error {
	_, err := db.Exec(`
		UPDATE communities 
		SET moderators = moderators || $1::jsonb
		WHERE community_id = $2 AND NOT moderators @> $3::jsonb`,
		fmt.Sprintf("[%d]", userID), communityID, fmt.Sprintf("[%d]", userID))
	return err
}
