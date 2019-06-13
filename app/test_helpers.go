package app

import (
	"errors"
	"log"
)

type TestServiceOptions struct {
	ForceAddQuestionError         bool
	ForceAddAnswerError           bool
	ForceCanUserAnswerQuestion    bool
	ForceCreateEventError         bool
	ForceCreateFeedbackError      bool
	ForceGetEventError            bool
	ForceGetEventsError           bool
	ForceGetQuestionError         bool
	ForceGetFeedbackByTokError    bool
	ForceGetFeedbackByTokNotFound bool
	ForceGetQuestionsForTokError  bool
	ForceMarkFeedbackAbsentError  bool
}

type TestDeliveryService struct {
	sendFeedbackEmailCount      int
	forceSendFeedbackEmailError bool
}

func (tds *TestDeliveryService) SendFeedbackEmail(string, string) error {
	tds.sendFeedbackEmailCount++
	log.Print("~~~sending feedback email")
	if tds.forceSendFeedbackEmailError {
		return errors.New("forced test error")
	}

	return nil
}

func (tds *TestDeliveryService) GetFeedbackEmailCount() int {
	return tds.sendFeedbackEmailCount
}

type TestFeedbackService struct {
	forceCreateError              bool
	forceGetFeedbackByTokError    bool
	forceGetFeedbackByTokNotFound bool
	forceMarkFeedbackAbsentError  bool
	forceGetQuestionsForTokError  bool
	Feedback                      Feedback
}

func (tfs *TestFeedbackService) CreateFeedback(*Feedback) error {
	if tfs.forceCreateError {
		return errors.New("forced test error")
	}

	return nil
}

func (tfs *TestFeedbackService) GetFeedbackByTok(string) (*Feedback, error) {
	if tfs.forceGetFeedbackByTokError {
		return nil, errors.New("forced test error")
	}

	if tfs.forceGetFeedbackByTokNotFound {
		return nil, nil
	}

	return &tfs.Feedback, nil
}

func (tfs *TestFeedbackService) MarkFeedbackAbsent(*Feedback) error {
	if tfs.forceMarkFeedbackAbsentError {
		return errors.New("forced test error")
	}

	return nil
}

func (tfs *TestFeedbackService) GetQuestionsForTok(tok string) *[]Question {
	if tfs.forceGetQuestionsForTokError {
		return nil
	}

	return &tfs.Questions
}

type TestEventService struct {
	forceGetEventError         bool
	forceGetEventsError        bool
	forceGetQuestionError      bool
	forceCreateEventError      bool
	forceCanUserAnswerQuestion bool
	forceAddQuestionError      bool
	forceAddAnswerError        bool
	Event                      Event
	Question                   Question
	Answer                     Answer
}

func (tes *TestEventService) GetEvent(ID int64) *Event {
	if tes.forceGetEventError {
		return nil
	}

	return &tes.Event
}

func (tes *TestEventService) GetEvents() (*[]Event, error) {
	if tes.forceGetEventsError {
		return nil, errors.New("forced test error")
	}

	return &[]Event{tes.Event}, nil
}

func (tes *TestEventService) CreateEvent(*Event) error {
	if tes.forceCreateEventError {
		return errors.New("forced test error")
	}

	return nil
}

func (tes *TestEventService) AddQuestion(*Question) error {
	if tes.forceAddQuestionError {
		return errors.New("forced test error")
	}

	return nil
}

func (tes *TestEventService) GetQuestion(ID int64) *Question {
	if tes.forceGetQuestionError {
		return nil
	}

	return &tes.Question
}

func (tes *TestEventService) CanUserAnswerQuestion(ID int64, tok string) bool {
	return tes.forceCanUserAnswerQuestion
}

func (tes *TestEventService) AddAnswer(*Answer, int64) error {
	if tes.forceAddAnswerError {
		return errors.New("forced test error")
	}

	return nil
}
