package main

import "time"

// Event represents a situation about which a user would like anonymous feedback.
type Event struct {
	ID                 int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time
	Title              string    `json:"title" binding:"required,max=100"`
	Description        string    `json:"description" binding:"required,max=1000"`
	Time               time.Time `json:"scheduled_time" binding:"required"`
	UserID             int64
	OrganizerQuestions []Question
	Feedback           []Feedback
}
