package app

type AnonDependencies struct {
	Events   EventService
	Feedback FeedbackService
	Delivery DeliveryService
}

func (ad *AnonDependencies) setEvents(e EventService) error {
	ad.Events = e
	return nil
}

func (ad *AnonDependencies) setFeedback(f FeedbackService) error {
	ad.Feedback = f
	return nil
}

func MockSearchDependencies(opts TestServiceOptions) *AnonDependencies {
	deps := &AnonDependencies{
		Events: &TestEventService{
			forceCreateEventError:      opts.ForceCreateEventError,
			forceGetEventError:         opts.ForceGetEventError,
			forceGetEventsError:        opts.ForceGetEventsError,
			forceGetQuestionError:      opts.ForceGetQuestionError,
			forceAddQuestionError:      opts.ForceAddQuestionError,
			forceCanUserAnswerQuestion: opts.ForceCanUserAnswerQuestion,
		},
		Feedback: &TestFeedbackService{
			forceCreateError:              opts.ForceCreateFeedbackError,
			forceGetFeedbackByTokError:    opts.ForceGetFeedbackByTokError,
			forceGetFeedbackByTokNotFound: opts.ForceGetFeedbackByTokNotFound,
			forceMarkFeedbackAbsentError:  opts.ForceMarkFeedbackAbsentError,
			forceGetQuestionsForTokError:  opts.ForceGetQuestionsForTokError,
		},
		Delivery: &TestDeliveryService{}}
	// deps.InitLogger()
	// deps.InitEs()

	return deps
}
