package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Create a new post
func CreatePost(db *sql.DB, communityID int, picture, caption string, author int) (int, error) {
	var postID int

	// First verify the author exists
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", author).Scan(&userExists)
	if err != nil {
		return 0, err
	}
	if !userExists {
		return 0, fmt.Errorf("author with ID %d does not exist", author)
	}

	err = db.QueryRow(`
		INSERT INTO posts (community_id, picture, caption, author, time_posted, likes) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING post_id`,
		communityID, picture, caption, author, time.Now(), "[]").Scan(&postID)

	if err == nil {
		// Add post to community's posts array
		_, err = db.Exec(`
			UPDATE communities 
			SET posts = posts || $1::jsonb
			WHERE community_id = $2`,
			fmt.Sprintf("[%d]", postID), communityID)
	}

	return postID, err
}

// Get post by ID
func GetPost(db *sql.DB, postID int) (*Post, error) {
	var post Post
	var likesJSON string

	err := db.QueryRow(`
		SELECT post_id, community_id, picture, caption, author, time_posted, likes
		FROM posts 
		WHERE post_id = $1`,
		postID).Scan(&post.PostID, &post.CommunityID, &post.Picture, &post.Caption,
		&post.Author, &post.TimePosted, &likesJSON)

	if err != nil {
		return nil, err
	}

	// Parse JSON array
	json.Unmarshal([]byte(likesJSON), &post.Likes)

	return &post, nil
}

// Like a post
func LikePost(db *sql.DB, postID, userID int) error {
	_, err := db.Exec(`
		UPDATE posts 
		SET likes = likes || $1::jsonb
		WHERE post_id = $2 AND NOT likes @> $3::jsonb`,
		fmt.Sprintf("[%d]", userID), postID, fmt.Sprintf("[%d]", userID))
	return err
}

// Unlike a post
func UnlikePost(db *sql.DB, postID, userID int) error {
	_, err := db.Exec(`
		UPDATE posts 
		SET likes = likes - $1
		WHERE post_id = $2`,
		userID, postID)
	return err
}

// Check if user has posted today in a community
func HasPostedToday(db *sql.DB, userID, communityID int) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM posts 
		WHERE author = $1 AND community_id = $2 AND DATE(time_posted) = CURRENT_DATE`,
		userID, communityID).Scan(&count)
	return count > 0, err
}
