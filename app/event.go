package app

import "time"

// Event represents a situation about which a user would like anonymous feedback.
type Event struct {
	ID                 int64
	Title              string     `json:"title" binding:"required,max=200"`
	Description        string     `json:"description" binding:"required,max=5000"`
	Time               time.Time  `json:"scheduled_time" binding:"required"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
	OrganizerQuestions []Question
	Feedback           []Feedback
}

type EventService interface {
	GetQuestion(int64) *Question
	GetEvent(int64) *Event
	GetEvents() (*[]Event, error)
	CanUserAnswerQuestion(int64, string) bool
	CreateEvent(*Event) error
	AddQuestion(*Question) error
	AddAnswer(*Answer, int64) error
}

// EventPostParams represents the information about an Event that a user can create.
type EventPostParams struct {
	Title       string    `json:"title" form:"title" binding:"required,max=200"`
	Description string    `json:"description" form:"description" binding:"required,max=5000"`
	Time        time.Time `json:"scheduled_time" form:"scheduled_time" binding:"required"`
	Audience    []string  `json:"audience" form:"audience" binding:"required,min=1"`
}

// EventPutParams represents the information about an Event that a user can update.
type EventPutParams struct {
	ID          int64     `json:"id" binding:"required"`
	Title       string    `json:"title"`
	Description string    `json:"description" binding:"max=1000"`
	Time        time.Time `json:"scheduled_time"`
	OrganizerID int64
}
