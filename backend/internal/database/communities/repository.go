package communities

import (
	"database/sql"
)

// Create a new community
func CreateCommunity(db *sql.DB, picture, description, postTime, defaultPrompt string) (int, error) {
	var communityID int
	err := db.QueryRow(`
		INSERT INTO communities (picture, description, members, moderators, posts, post_time, default_prompt) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING community_id`,
		picture, description, "{}", "{}", "{}", postTime, defaultPrompt).Scan(&communityID)
	return communityID, err
}

// Get community by ID
func GetCommunity(db *sql.DB, communityID int) (*Community, error) {
	var community Community

	err := db.QueryRow(`
		SELECT community_id, picture, description, members, moderators, posts, post_time, default_prompt
		FROM communities 
		WHERE community_id = $1`,
		communityID).Scan(&community.CommunityID, &community.Picture, &community.Description,
		&community.Members, &community.Moderators, &community.Posts, &community.PostTime, &community.DefaultPrompt)

	if err != nil {
		return nil, err
	}

	return &community, nil
}

// Join a community
func JoinCommunity(db *sql.DB, userID, communityID int) error {
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
func LeaveCommunity(db *sql.DB, userID, communityID int) error {
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

// Add moderator to community
func AddModerator(db *sql.DB, communityID, userID int) error {
	_, err := db.Exec(`
		UPDATE communities 
		SET moderators = array_append(moderators, $1)
		WHERE community_id = $2 AND NOT $1 = ANY(moderators)`,
		userID, communityID)
	return err
}
