package dailies

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DEBUG: changed from community_id -> communityId... see if that breaks anything
type Daily struct {
	PostID      string    `json:"posd" db:"post_id"`
	CommunityID string    `json:"communityId" db:"community_id"`
	Picture     string    `json:"picture"`
	Caption     string    `json:"caption"`
	Author      string    `json:"author"`
	TimePosted  time.Time `json:"timePosted" db:"time_posted"`
	Likes       []string  `json:"likes"`
}

// Fetch all dailies from certain community
func FetchDailiesFromCommunity(db *pgxpool.Pool, communityID string) ([]Daily, error) {
	rows, err := db.Query(
		context.TODO(),
		`SELECT * FROM dailies WHERE community_id = $1`,
		communityID)
	if err != nil {
		return nil, err
	}
	dailies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Daily])
	if err != nil {
		return nil, err
	}
	return dailies, nil
}

// Fetch all dailies from certain user
func FetchDailiesFromUser(db *pgxpool.Pool, userId string) ([]Daily, error) {
	rows, err := db.Query(
		context.TODO(),
		`SELECT * FROM dailies WHERE author = $1`,
		userId)
	if err != nil {
		return nil, err
	}
	dailies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Daily])
	if err != nil {
		return nil, err
	}
	return dailies, nil
}

// TODO: should this check if daily already exists for that user
// Create a new daily
func CreateDaily(db *pgxpool.Pool, communityID string, pictureID, caption string, author string) (string, error) {
	var postID string

	// Verify the author exists
	var userExists bool
	err := db.QueryRow(context.TODO(), "SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", author).Scan(&userExists)
	if err != nil {
		return "", err
	}
	if !userExists {
		return "", fmt.Errorf("author with ID %s does not exist", author)
	}

	// Generate a unique post ID
	postID = uuid.New().String()

	// Insert into dailies table (picture now stores the file path)
	_, err = db.Exec(context.TODO(), `
		INSERT INTO dailies (post_id, community_id, picture, caption, author, time_posted, likes) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		postID, communityID, pictureID, caption, author, time.Now(), "{}")
	if err == nil {
		// Add daily to community's posts array
		_, err = db.Exec(context.TODO(), `
			UPDATE communities 
			SET posts = array_append(posts, $1)
			WHERE community_id = $2`,
			postID, communityID)
	}
	return postID, err
}

// Get daily by ID
func GetDaily(db *pgxpool.Pool, postID string) (*Daily, error) {
	var daily Daily

	err := db.QueryRow(context.TODO(), `
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
func LikeDaily(db *pgxpool.Pool, postID, userID string) error {
	// Add user to likes array if not already present
	_, err := db.Exec(context.TODO(), `
		UPDATE dailies 
		SET likes = array_append(likes, $1)
		WHERE post_id = $2 AND NOT $1 = ANY(likes)`,
		userID, postID)
	return err
}

// Unlike a daily
func UnlikeDaily(db *pgxpool.Pool, postID, userID string) error {
	// Remove user from likes array
	_, err := db.Exec(context.TODO(), `
		UPDATE dailies 
		SET likes = array_remove(likes, $1)
		WHERE post_id = $2`,
		userID, postID)
	return err
}

// Check if user has posted today in a community
func HasPostedToday(db *pgxpool.Pool, userID, communityID string) (bool, error) {
	var count int
	err := db.QueryRow(context.TODO(), `
		SELECT COUNT(*) FROM dailies 
		WHERE author = $1 AND community_id = $2`,
		userID, communityID).Scan(&count)
	return count > 0, err
}

// TODO: no confirmation of deletion
func DeleteDailyByUser(db *pgxpool.Pool, userId, communityId string) error {
	// delete from daily db
	var postId string
	err := db.QueryRow(context.TODO(), `
		DELETE FROM dailies
		WHERE author = $1
		AND community_id = $2
		RETURNING post_id`,
		userId, communityId).Scan(&postId)
	if err != nil {
		return err
	}
	// remove from community posts
	_, err = db.Exec(context.TODO(), `
		UPDATE communities
		SET posts = array_remove(posts, $1)
		WHERE community_id = $2`,
		postId, communityId)
	return err
}

func DeleteDaily(db *pgxpool.Pool, postID string) error {
	tx, err := db.Begin(context.TODO())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.TODO())

	// First get the community ID to remove the post from community's posts array
	var communityID string
	err = tx.QueryRow(context.TODO(), `
		SELECT community_id FROM dailies WHERE post_id = $1`,
		postID).Scan(&communityID)
	if err != nil {
		return fmt.Errorf("failed to find daily: %v", err)
	}

	// Remove the post from the community's posts array
	_, err = tx.Exec(context.TODO(), `
		UPDATE communities 
		SET posts = array_remove(posts, $1)
		WHERE community_id = $2`,
		postID, communityID)
	if err != nil {
		return fmt.Errorf("failed to remove post from community: %v", err)
	}

	// Delete the daily post
	tag, err := tx.Exec(context.TODO(), `
		DELETE FROM dailies 
		WHERE post_id = $1`,
		postID)
	if err != nil {
		return fmt.Errorf("failed to delete daily: %v", err)
	}

	// Check if any row was actually deleted
	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("daily post with ID %s not found", postID)
	}

	return tx.Commit(context.TODO())
}
