// Connecting only worked when I put this in terminal:
// $env:PG_DSN = "postgresql://SQLDatabase:BinghamtonACMPT4@database-1.c5qccgo64zmj.us-east-2.rds.amazonaws.com:5432/postgres?sslmode=require"
// Also only worked when I added a rule to allow any IP to connect to the AWS.
// This is mostly AI and I am still struggling to understand the majority of this, but it works.

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ================ DATABASE STRUCTS ===========================
type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profile_picture"`
	JoinedDate     time.Time `json:"joined_date"`
	Communities    []int     `json:"communities"` // Array of community IDs
	Friends        []int     `json:"friends"`     // Array of user IDs
}

// Post struct represents a daily post
type Post struct {
	PostID      int       `json:"post_id"`
	CommunityID int       `json:"community_id"`
	Picture     string    `json:"picture"`
	Caption     string    `json:"caption"`
	Author      int       `json:"author"` // User ID
	TimePosted  time.Time `json:"time_posted"`
	Likes       []int     `json:"likes"` // Array of user IDs who liked the post
}

// Community struct represents a community
type Community struct {
	CommunityID   int    `json:"community_id"`
	Picture       string `json:"picture"`
	Description   string `json:"description"`
	Members       []int  `json:"members"`    // Array of user IDs
	Moderators    []int  `json:"moderators"` // Array of user IDs
	Posts         []int  `json:"posts"`      // Array of post IDs
	PostTime      string `json:"post_time"`  // Default time for daily posts (e.g., "09:00:00")
	DefaultPrompt string `json:"default_prompt"`
}

// ==================== DATABASE VERIFICATION ====================

func VerifyDatabaseStructure(db *sql.DB) error {
	requiredTables := []string{"communities", "posts", "users"}

	for _, table := range requiredTables {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' s
				AND table_name = $1
			)`, table).Scan(&exists)

		if err != nil {
			return fmt.Errorf("failed to check table %s: %v", table, err)
		}

		if !exists {
			return fmt.Errorf("required table '%s' does not exist. Run setup_database.go first", table)
		}
		fmt.Printf("✅ Table '%s' exists\n", table)
	}

	return nil
}

func ListAllTables(db *sql.DB) error {
	fmt.Println("\n📊 Listing all tables in the database:")
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
		ORDER BY table_name;
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tableName string
	tableCount := 0
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		fmt.Printf(" - %s\n", tableName)
		tableCount++
	}

	if tableCount == 0 {
		fmt.Println("No tables found in the database.")
	} else {
		fmt.Printf("Total tables: %d\n", tableCount)
	}
	return nil
}

func ShowTableStructure(db *sql.DB, tableName string) error {
	fmt.Printf("\n📋 Structure of '%s' table:\n", tableName)
	columns, err := db.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns 
		WHERE table_name = $1 AND table_schema = 'public'
		ORDER BY ordinal_position;
	`, tableName)
	if err != nil {
		return err
	}
	defer columns.Close()

	var colName, dataType, nullable string
	fmt.Println("Column Name     | Data Type     | Nullable")
	fmt.Println("----------------|---------------|----------")
	for columns.Next() {
		if err := columns.Scan(&colName, &dataType, &nullable); err != nil {
			return err
		}
		fmt.Printf("%-15s | %-13s | %-8s\n", colName, dataType, nullable)
	}
	return nil
}

// ==================== USER OPERATIONS ====================

// Register a new user
func Register(db *sql.DB, name, email, profilePicture string) (int, error) {
	var userID int

	// Check if user already exists
	var existingID int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		// User already exists, return the existing ID
		return existingID, nil
	} else if err != sql.ErrNoRows {
		// Some other error occurred
		return 0, err
	}

	// User doesn't exist, create new one
	err = db.QueryRow(`
        INSERT INTO users (name, email, profile_picture, joined_date, communities, friends, created_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id`,
		name, email, profilePicture, time.Now(), "[]", "[]", time.Now()).Scan(&userID)
	return userID, err
}

// Login user by email
func Login(db *sql.DB, email string) (*User, error) {
	var user User
	var communitiesJSON, friendsJSON string

	err := db.QueryRow(`
        SELECT id, name, email, profile_picture, joined_date, communities, friends
        FROM users 
        WHERE email = $1`,
		email).Scan(&user.ID, &user.Name, &user.Email, &user.ProfilePicture,
		&user.JoinedDate, &communitiesJSON, &friendsJSON)

	if err != nil {
		return nil, err
	}

	// Parse JSON arrays
	json.Unmarshal([]byte(communitiesJSON), &user.Communities)
	json.Unmarshal([]byte(friendsJSON), &user.Friends)

	return &user, nil
}

// Follow another user
func Follow(db *sql.DB, userID, friendID int) error {
	_, err := db.Exec(`
		UPDATE users 
		SET friends = friends || $1::jsonb
		WHERE id = $2 AND NOT friends @> $3::jsonb`,
		fmt.Sprintf("[%d]", friendID), userID, fmt.Sprintf("[%d]", friendID))
	return err
}

// Unfollow a user
func Unfollow(db *sql.DB, userID, friendID int) error {
	_, err := db.Exec(`
		UPDATE users 
		SET friends = friends - $1
		WHERE id = $2`,
		friendID, userID)
	return err
}

// Change profile picture
func ChangeProfile(db *sql.DB, userID int, profilePicture string) error {
	_, err := db.Exec("UPDATE users SET profile_picture = $1 WHERE id = $2", profilePicture, userID)
	return err
}

// Delete user account
func DeleteAccount(db *sql.DB, userID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove user from all communities
	_, err = tx.Exec(`UPDATE communities SET members = members - $1`, userID)
	if err != nil {
		return err
	}

	// Remove user as moderator from all communities
	_, err = tx.Exec(`UPDATE communities SET moderators = moderators - $1`, userID)
	if err != nil {
		return err
	}

	// Delete user's posts
	_, err = tx.Exec("DELETE FROM posts WHERE author = $1", userID)
	if err != nil {
		return err
	}

	// Remove user from friends lists
	_, err = tx.Exec(`UPDATE users SET friends = friends - $1`, userID)
	if err != nil {
		return err
	}

	// Delete the user
	_, err = tx.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ==================== COMMUNITY OPERATIONS ====================

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

// ==================== POST OPERATIONS ====================

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
