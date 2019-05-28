package controllers

import "testing"

func TestPostQuestion(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		//TODO
	})
	t.Run("NoTokenPayload", func(t *testing.T) {

	})
	t.Run("UserIsNotOwner", func(t *testing.T) {

	})
	t.Run("EventIdsDoNotMatch", func(t *testing.T) {

	})
	t.Run("EventNotFound", func(t *testing.T) {

	})
	t.Run("ValidationErrors", func(t *testing.T) {
		t.Run("NoTitleProvided", func(t *testing.T) {

		})
		t.Run("TitleTooLong", func(t *testing.T) {

		})
		t.Run("ContentTooLong", func(t *testing.T) {

		})
		t.Run("AnswersTooLong", func(t *testing.T) {

		})
	})
	t.Run("ErrorSavingQuestion", func(t *testing.T) {

	})
	t.Run("Success", func(t *testing.T) {

	})
}
