package users

import (
	"database/sql"
	"fmt"
	"time"
)

func GetUser(db *sql.DB, userID string) (*User, error) {
	var user User

	err := db.QueryRow(`
		SELECT user_id, name, email, profile_picture, joined_date, communities
		FROM users 
		WHERE user_id = $1`,
		userID).Scan(&user.UserID, &user.Name, &user.Email, &user.ProfilePicture,
		&user.JoinedDate, &user.Communities)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Register a new user
func Register(db *sql.DB, name, email, profilePicture string) (string, error) {
	// Check if user already exists
	var existingID string
	err := db.QueryRow("SELECT user_id FROM users WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		// User already exists, return the existing ID
		return existingID, nil
	} else if err != sql.ErrNoRows {
		// Some other error occurred
		return "", err
	}

	// Generate a unique user ID
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())

	// User doesn't exist, create new one
	_, err = db.Exec(`
        INSERT INTO users (user_id, name, email, profile_picture, joined_date, communities) 
        VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, name, email, profilePicture, time.Now(), "{}")
	return userID, err
}

// Login user by email
func Login(db *sql.DB, email string) (*User, error) {
	var user User

	err := db.QueryRow(`
        SELECT user_id, name, email, profile_picture, joined_date, communities
        FROM users 
        WHERE email = $1`,
		email).Scan(&user.UserID, &user.Name, &user.Email, &user.ProfilePicture,
		&user.JoinedDate, &user.Communities)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Change profile picture
func ChangeProfile(db *sql.DB, userID string, profilePicture string) error {
	_, err := db.Exec("UPDATE users SET profile_picture = $1 WHERE user_id = $2", profilePicture, userID)
	return err
}

// Delete user account
func DeleteAccount(db *sql.DB, userID string) error {
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

	// Delete user's dailies
	_, err = tx.Exec("DELETE FROM dailies WHERE author = $1", userID)
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
