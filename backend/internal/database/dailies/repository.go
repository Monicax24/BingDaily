package dailies

import (
	"database/sql"
	"fmt"
	"time"
)

// Create a new daily
func CreateDaily(db *sql.DB, communityID int, picture, caption string, author int) (int, error) {
	var postID int

	// First verify the author exists
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", author).Scan(&userExists)
	if err != nil {
		return 0, err
	}
	if !userExists {
		return 0, fmt.Errorf("author with ID %d does not exist", author)
	}

	// Insert into dailies table
	err = db.QueryRow(`
		INSERT INTO dailies (community_id, picture, caption, author, time_posted, likes) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING post_id`,
		communityID, picture, caption, author, time.Now(), 0).Scan(&postID)

	if err == nil {
		// Add daily to community's posts array using PostgreSQL array append
		_, err = db.Exec(`
			UPDATE communities 
			SET posts = array_append(posts, $1)
			WHERE community_id = $2`,
			postID, communityID)
	}

	return postID, err
}

// Get daily by ID
func GetDaily(db *sql.DB, postID int) (*Daily, error) {
	var daily Daily

	err := db.QueryRow(`
		SELECT post_id, community_id, picture, caption, author, time_posted, likes
		FROM dailies 
		WHERE post_id = $1`,
		postID).Scan(&daily.PostID, &daily.CommunityID, &daily.Picture, &daily.Caption,
		&daily.Author, &daily.TimePosted, &daily.Likes)

	if err != nil {
		return nil, err
	}

	return &daily, nil
}

// Like a daily
func LikeDaily(db *sql.DB, postID, userID int) error {
	// Increment like count
	_, err := db.Exec(`
		UPDATE dailies 
		SET likes = likes + 1
		WHERE post_id = $1`,
		postID)
	return err
}

// Unlike a daily
func UnlikeDaily(db *sql.DB, postID, userID int) error {
	// Decrement like count (ensure it doesn't go below 0)
	_, err := db.Exec(`
		UPDATE dailies 
		SET likes = GREATEST(likes - 1, 0)
		WHERE post_id = $1`,
		postID)
	return err
}

// Check if user has posted today in a community
func HasPostedToday(db *sql.DB, userID, communityID int) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM dailies 
		WHERE author = $1 AND community_id = $2 AND DATE(time_posted) = CURRENT_DATE`,
		userID, communityID).Scan(&count)
	return count > 0, err
}
