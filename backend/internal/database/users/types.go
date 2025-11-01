package users

import "time"

type User struct {
	UserID         int       `json:"user_id"` // Changed from ID string to UserID int
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profile_picture"`
	JoinedDate     time.Time `json:"joined_date"`
	Communities    []int     `json:"communities"` // Array of community IDs
	Friends        []int     `json:"friends"`     // Array of user IDs
}
