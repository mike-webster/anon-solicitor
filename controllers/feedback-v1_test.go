package controllers

import "testing"

func TestGetFeedbackV1(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		//TODO
	})

	t.Run("TokenError", func(t *testing.T) {
		// TODO
	})

	t.Run("IdLessThanOne", func(t *testing.T) {
		// TODO
	})

	t.Run("EventNotFound", func(t *testing.T) {
		// TODO
	})

	t.Run("GetFeedbackByTokError", func(t *testing.T) {
		// TODO
	})

	t.Run("FeedbackNotFound", func(t *testing.T) {
		// TODO
	})

	t.Run("Success", func(t *testing.T) {
		// TODO
	})
}

func TestPostAbsentFeedbackV1(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		// TODO
	})

	t.Run("TokenError", func(t *testing.T) {
		// TODO
	})

	t.Run("BlankToken", func(t *testing.T) {
		// TODO
	})

	t.Run("GetFeedbackByTokError", func(t *testing.T) {
		// TODO
	})

	t.Run("MarkFeedbackAbsentError", func(t *testing.T) {
		// TODO
	})

	t.Run("Success", func(t *testing.T) {
		// TODO
	})
}

func TestPostFeedbackV1(t *testing.T) {
	// TODO
}
