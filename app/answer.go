package app

import "time"

type Answer struct {
	ID         int64
	QuestionID int64
	Content    string
	CreatedAt  time.Time
}

type AnswerService interface {
	AnswerQuestion(*Answer) error
}

type AnswerPostParams struct {
	Content string `json:"content" binding:"required,max=5000"`
}
