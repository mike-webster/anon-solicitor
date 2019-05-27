package app

import "time"

// Question fd
type Question struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type QuestionService interface {
	Create(*Question) error
}
