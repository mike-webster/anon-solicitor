package app

import "time"

type Answer struct {
	ID         int64     `db:"id"`
	QuestionID int64     `db:"question_id"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
}

type AnswerService interface {
	AnswerQuestion(*Answer) error
}

type AnswerPostParams struct {
	Content string `json:"content" binding:"required,max=5000"`
}
