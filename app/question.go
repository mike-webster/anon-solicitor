package app

import "time"

// Question fd
type Question struct {
	ID        int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type QuestionService interface {
	CreateQuestion(*Question) error
}
