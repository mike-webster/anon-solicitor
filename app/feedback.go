package app

// Feedback represents a user's opinion on an event.
type Feedback struct {
	ID      int64  `db:"id"`
	Tok     string `db:"tok"`
	EventID int64  `db:"event_id"`
	Content string `db:"content"`
	Absent  bool   `db:"absent"`
}

type FeedbackService interface {
	CreateFeedback(*Feedback) error
	GetFeedbackByTok(string) (*Feedback, error)
	MarkFeedbackAbsent(*Feedback) error
}
