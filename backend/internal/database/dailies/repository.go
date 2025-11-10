package dailies

import (
	"database/sql"
	"fmt"
	"time"
)

// Create a new daily
func CreateDaily(db *sql.DB, communityID string, picturePath, caption string, author string) (string, error) {
	var postID string

	// Verify the author exists
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", author).Scan(&userExists)
	if err != nil {
		return "", err
	}
	if !userExists {
		return "", fmt.Errorf("author with ID %s does not exist", author)
	}

	// Generate a unique post ID
	postID = fmt.Sprintf("post_%d", time.Now().UnixNano())

	// Insert into dailies table (picture now stores the file path)
	_, err = db.Exec(`
		INSERT INTO dailies (post_id, community_id, picture, caption, author, time_posted, likes) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		postID, communityID, picturePath, caption, author, time.Now(), "{}")

	if err == nil {
		// Add daily to community's posts array
		_, err = db.Exec(`
			UPDATE communities 
			SET posts = array_append(posts, $1)
			WHERE community_id = $2`,
			postID, communityID)
	}

	return postID, err
}

// Get daily by ID
func GetDaily(db *sql.DB, postID string) (*Daily, error) {
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
func LikeDaily(db *sql.DB, postID, userID string) error {
	// Add user to likes array if not already present
	_, err := db.Exec(`
		UPDATE dailies 
		SET likes = array_append(likes, $1)
		WHERE post_id = $2 AND NOT $1 = ANY(likes)`,
		userID, postID)
	return err
}

// Unlike a daily
func UnlikeDaily(db *sql.DB, postID, userID string) error {
	// Remove user from likes array
	_, err := db.Exec(`
		UPDATE dailies 
		SET likes = array_remove(likes, $1)
		WHERE post_id = $2`,
		userID, postID)
	return err
}

// Check if user has posted today in a community
func HasPostedToday(db *sql.DB, userID, communityID string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM dailies 
		WHERE author = $1 AND community_id = $2 AND DATE(time_posted) = CURRENT_DATE`,
		userID, communityID).Scan(&count)
	return count > 0, err
}
func DeleteDaily(db *sql.DB, postID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First get the community ID to remove the post from community's posts array
	var communityID string
	err = tx.QueryRow(`
		SELECT community_id FROM dailies WHERE post_id = $1`,
		postID).Scan(&communityID)
	if err != nil {
		return fmt.Errorf("failed to find daily: %v", err)
	}

	// Remove the post from the community's posts array
	_, err = tx.Exec(`
		UPDATE communities 
		SET posts = array_remove(posts, $1)
		WHERE community_id = $2`,
		postID, communityID)
	if err != nil {
		return fmt.Errorf("failed to remove post from community: %v", err)
	}

	// Delete the daily post
	result, err := tx.Exec(`
		DELETE FROM dailies 
		WHERE post_id = $1`,
		postID)
	if err != nil {
		return fmt.Errorf("failed to delete daily: %v", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("daily post with ID %s not found", postID)
	}

	return tx.Commit()
}
