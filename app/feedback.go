package app

// Feedback represents a user's opinion on an event.
type Feedback struct {
	ID      int64
	Tok     string
	EventID int64
}

type FeedbackService interface {
	CreateFeedback(*Feedback) error
	GetFeedbackByTok(string) (*Feedback, error)
	MarkFeedbackAbsent(*Feedback) error
}
