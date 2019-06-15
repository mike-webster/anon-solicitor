package app

import (
	"database/sql"
	"strings"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

// Question fd
type Question struct {
	ID        int64
	EventID   int64  `db:"event_id"`
	Title     string `db:"title"`
	Content   string
	DBAnswers sql.NullString `db:"answers"` // This is for "options" for the question, delimited by ";;"
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}

type QuestionDisplay struct {
	QuestionID int64
	EventID    int64 `form:"eventid"`
	Content    string
	Answers    []string
	Feedback   string
}

func (q *Question) AddAnswers(answers string) {
	q.DBAnswers = sql.NullString{String: answers}
}

func (q *Question) Answers() []string {
	if q.DBAnswers.Valid {
		return strings.Split(q.DBAnswers.String, ";")
	}
	return []string{}
}

type QuestionPostParams struct {
	Title   string `json:"title" binding:"required,max=5000"`
	Content string `json:"content" binding:"max=5000"`
	Answers string `json:"answers" binding:"max=5000"` // This is for "options" for the question, delimited by ";;"
}

type QuestionService interface {
	Create(*Question) error
}
