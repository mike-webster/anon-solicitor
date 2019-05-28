package app

import (
	"errors"
	"log"
)

type TestServiceOptions struct {
	ForceGetEventError            bool
	ForceGetEventsError           bool
	ForceCreateEventError         bool
	ForceCreateFeedbackError      bool
	ForceGetFeedbackByTokError    bool
	ForceMarkFeedbackAbsentError  bool
	ForceGetFeedbackByTokNotFound bool
}

type TestEventService struct {
	forceGetEventError    bool
	forceGetEventsError   bool
	forceCreateEventError bool
	forceAddQuestionError bool
	Event                 Event
}

func (tes *TestEventService) setNoErrors() {
	tes.forceGetEventError = false
	tes.forceCreateEventError = false
	tes.forceGetEventsError = false
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
