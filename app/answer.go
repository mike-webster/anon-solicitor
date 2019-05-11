package app

type Answer struct {
	QuestionID int64
	EventID    int64
	Content    string
}

type AnswerService interface {
	AnswerQuestion(*Answer) error
}
