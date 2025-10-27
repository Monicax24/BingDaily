package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

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
