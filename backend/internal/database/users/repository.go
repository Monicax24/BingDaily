package users

import (
	"database/sql"
	"time"
)

// Register a new user
func Register(db *sql.DB, name, email, profilePicture string) (int, error) {
	var userID int

	// Check if user already exists
	var existingID int
	err := db.QueryRow("SELECT user_id FROM users WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		// User already exists, return the existing ID
		return existingID, nil
	} else if err != sql.ErrNoRows {
		// Some other error occurred
		return 0, err
	}

	// User doesn't exist, create new one
	err = db.QueryRow(`
        INSERT INTO users (name, email, profile_picture, joined_date, communities, friends) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING user_id`,
		name, email, profilePicture, time.Now(), "{}", "{}").Scan(&userID)
	return userID, err
}

// Login user by email
func Login(db *sql.DB, email string) (*User, error) {
	var user User

	err := db.QueryRow(`
        SELECT user_id, name, email, profile_picture, joined_date, communities, friends
        FROM users 
        WHERE email = $1`,
		email).Scan(&user.UserID, &user.Name, &user.Email, &user.ProfilePicture,
		&user.JoinedDate, &user.Communities, &user.Friends)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Follow another user
func Follow(db *sql.DB, userID, friendID int) error {
	_, err := db.Exec(`
		UPDATE users 
		SET friends = array_append(friends, $1)
		WHERE user_id = $2 AND NOT $1 = ANY(friends)`,
		friendID, userID)
	return err
}

// Unfollow a user
func Unfollow(db *sql.DB, userID, friendID int) error {
	_, err := db.Exec(`
		UPDATE users 
		SET friends = array_remove(friends, $1)
		WHERE user_id = $2`,
		friendID, userID)
	return err
}

// Change profile picture
func ChangeProfile(db *sql.DB, userID int, profilePicture string) error {
	_, err := db.Exec("UPDATE users SET profile_picture = $1 WHERE user_id = $2", profilePicture, userID)
	return err
}

// Delete user account
func DeleteAccount(db *sql.DB, userID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove user from all communities' members
	_, err = tx.Exec(`UPDATE communities SET members = array_remove(members, $1)`, userID)
	if err != nil {
		return err
	}

	// Remove user from all communities' moderators
	_, err = tx.Exec(`UPDATE communities SET moderators = array_remove(moderators, $1)`, userID)
	if err != nil {
		return err
	}

	// Delete user's dailies
	_, err = tx.Exec("DELETE FROM dailies WHERE author = $1", userID)
	if err != nil {
		return err
	}

	// Remove user from friends lists
	_, err = tx.Exec(`UPDATE users SET friends = array_remove(friends, $1)`, userID)
	if err != nil {
		return err
	}

	// Delete the user
	_, err = tx.Exec("DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
