package users

import "time"

type User struct {
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profile_picture"`
	JoinedDate     time.Time `json:"joined_date"`
	Communities    []string  `json:"communities"` // Array of community IDs
}
