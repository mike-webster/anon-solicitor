package models

import "time"

// User represents an application user.
type User struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Email     string
	Active    bool
}
