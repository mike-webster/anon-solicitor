package app

import "time"

// Question fd
type Question struct {
	ID        int64
	EventID   int64
	Title     string
	Content   string
	Answers   string // This is for "options" for the question, delimited by ";;"
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type QuestionPostParams struct {
	Title   string `json:"title" binding:"required,max=5000"`
	Content string `json:"content" binding:"max=5000"`
	Answers string `json:"answers" binding:"max=5000"` // This is for "options" for the question, delimited by ";;"
}

type QuestionService interface {
	Create(*Question) error
}
