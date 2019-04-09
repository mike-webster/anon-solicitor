package models

// Feedback represents a user's opinion on an event.
type Feedback struct {
	ID      int64
	Content string
	EventID int64
}
