package app

import (
	"errors"
)

type TestServiceOptions struct {
	ForceGetEventError    bool
	ForceGetEventsError   bool
	ForceCreateEventError bool
}

type TestEventService struct {
	forceGetEventError    bool
	forceGetEventsError   bool
	forceCreateEventError bool
	Event                 Event
}

func (tes *TestEventService) setNoErrors() {
	tes.forceGetEventError = false
	tes.forceCreateEventError = false
	tes.forceGetEventsError = false
}

type TestFeedbackService struct {
	forceCreateError           bool
	forceGetFeedbackByTokError bool
	forceMarkFeedbackAbsent    bool
	Feedback                   Feedback
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

	return &tfs.Feedback, nil
}

func (tfs *TestFeedbackService) MarkFeedbackAbsent(*Feedback) error {
	if tfs.forceMarkFeedbackAbsent {
		return errors.New("forced test error")
	}

	return nil
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
