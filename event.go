package anon

import "time"

// Event represents a situation about which a user would like anonymous feedback.
type Event struct {
	ID                 int64
	Title              string     `json:"title" binding:"required,max=200"`
	Description        string     `json:"description" binding:"required,max=5000"`
	Time               time.Time  `json:"scheduled_time" binding:"required"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
	UserID             int64      `json:"user_id" db:"user_id"`
	OrganizerQuestions []Question
	Feedback           []Feedback
}

type EventService interface {
	GetEvent(id int64) *Event
	CreateEvent(event *Event) error
	UpdateEvent(event *Event) error
}

// EventPostParams represents the information about an Event that a user can create.
type EventPostParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description" binding:"max=1000"`
	Time        time.Time `json:"scheduled_time"`
}

// EventPutParams represents the information about an Event that a user can update.
type EventPutParams struct {
	ID          int64     `json:"id" binding:"required"`
	Title       string    `json:"title"`
	Description string    `json:"description" binding:"max=1000"`
	Time        time.Time `json:"scheduled_time"`
	OrganizerID int64
}
