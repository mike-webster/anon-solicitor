package controllers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/mike-webster/anon-solicitor/app"
)

func TestPostQuestion(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		//TODO
	})
	t.Run("NoTokenPayload", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		headers["token"] = getTestTok(nil)
		req := performRequest(r, "POST", "/v1/questions/1", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})
	t.Run("UserIsNotOwner", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleAudience})
		req := performRequest(r, "POST", "/v1/questions/1", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})
	t.Run("EventIdsDoNotMatch", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 2})
		req := performRequest(r, "POST", "/v1/questions/1", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})
	t.Run("EventNotFound", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceGetEventError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000, "qid": 1})
		req := performRequest(r, "POST", "/v1/answers/10000/1", nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})
	t.Run("ValidationErrors", func(t *testing.T) {
		t.Run("NoTitleProvided", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{}
			question := app.QuestionPostParams{
				Content: "test content",
			}
			bytes, _ := json.Marshal(question)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000})
			req := performRequest(r, "POST", "/v1/questions/10000", &bytes, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body)
		})
		t.Run("TitleTooLong", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{}
			var longTitle string
			for i := 0; i < 5001; i++ {
				longTitle += "a"
			}
			question := app.QuestionPostParams{
				Title:   longTitle,
				Content: "test content",
			}
			bytes, _ := json.Marshal(question)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000})
			req := performRequest(r, "POST", "/v1/questions/10000", &bytes, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body)
		})
		t.Run("ContentTooLong", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{}
			var longContent string
			for i := 0; i < 5001; i++ {
				longContent += "a"
			}
			question := app.QuestionPostParams{
				Title:   "test title",
				Content: longContent,
			}
			bytes, _ := json.Marshal(question)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000})
			req := performRequest(r, "POST", "/v1/questions/10000", &bytes, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body)
		})
		t.Run("AnswersTooLong", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{}
			var longAnswers string
			for i := 0; i < 5001; i++ {
				longAnswers += "a"
			}
			question := app.QuestionPostParams{
				Title:   "test title",
				Content: "test content",
				Answers: longAnswers,
			}
			bytes, _ := json.Marshal(question)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000})
			req := performRequest(r, "POST", "/v1/questions/10000", &bytes, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body)
		})
	})
	t.Run("ErrorSavingQuestion", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceAddQuestionError: true,
		}
		question := app.QuestionPostParams{
			Title:   "test title",
			Content: "test content",
		}
		bytes, _ := json.Marshal(question)
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000})
		req := performRequest(r, "POST", "/v1/questions/10000", &bytes, headers)
		assert.Equal(t, http.StatusInternalServerError, req.Code, req.Body)
	})
	t.Run("Success", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		question := app.QuestionPostParams{
			Title:   "test title",
			Content: "test content",
		}
		bytes, _ := json.Marshal(question)
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000})
		req := performRequest(r, "POST", "/v1/questions/10000", &bytes, headers)
		assert.Equal(t, http.StatusOK, req.Code, req.Body)
	})
}

func TestPostAnswer(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		//TODO
	})
	t.Run("NoToken", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		headers["token"] = getTestTok(nil)
		req := performRequest(r, "POST", "/v1/answers/1/1", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})
	t.Run("RoleNotPermitted", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		headers["token"] = getTestTok(&map[string]interface{}{"role": "notARealRole"})
		req := performRequest(r, "POST", "/v1/answers/1/1", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})
	t.Run("MismatchedEventIDs", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 2})
		req := performRequest(r, "POST", "/v1/answers/1/1", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})
	t.Run("EventNotFound", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceGetEventError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 10000})
		req := performRequest(r, "POST", "/v1/answers/10000/1", nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})
	t.Run("MismatchedQuestionIDs", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 1, "qid": 2})
		req := performRequest(r, "POST", "/v1/answers/1/1", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})
	t.Run("QuestionNotFound", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceGetQuestionError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 1, "qid": 1})
		req := performRequest(r, "POST", "/v1/answers/1/1", nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})
	t.Run("UserCantAnswerQuestion", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceCanUserAnswerQuestion: false,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 1, "qid": 1})
		req := performRequest(r, "POST", "/v1/answers/1/1", nil, headers)
		assert.Equal(t, http.StatusForbidden, req.Code)
	})
	t.Run("ValidationErrors", func(t *testing.T) {
		t.Run("NoContent", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{
				ForceCanUserAnswerQuestion: true,
			}
			answer := app.AnswerPostParams{}
			bytes, _ := json.Marshal(answer)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 1, "qid": 1})
			req := performRequest(r, "POST", "/v1/answers/1/1", &bytes, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code)
		})
		t.Run("ContentTooLong", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{
				ForceCanUserAnswerQuestion: true,
			}
			answer := app.AnswerPostParams{}
			for i := 0; i < 5001; i++ {
				answer.Content += "a"
			}
			bytes, _ := json.Marshal(answer)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 1, "qid": 1})
			req := performRequest(r, "POST", "/v1/answers/1/1", &bytes, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code)
		})
	})
	t.Run("CantSaveAnswer", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceCanUserAnswerQuestion: true,
			ForceAddAnswerError:        true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 1, "qid": 1})
		req := performRequest(r, "POST", "/v1/answers/1/1", nil, headers)
		assert.Equal(t, http.StatusInternalServerError, req.Code)
	})
	t.Run("Success", func(t *testing.T) {
		t.Run("Owner", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{
				ForceCanUserAnswerQuestion: true,
			}
			answer := app.AnswerPostParams{
				Content: "test content",
			}
			bytes, _ := json.Marshal(answer)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleOwner, "eid": 1, "qid": 1})
			req := performRequest(r, "POST", "/v1/answers/1/1", &bytes, headers)
			assert.Equal(t, http.StatusOK, req.Code)
		})
		t.Run("Audience", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{
				ForceCanUserAnswerQuestion: true,
			}
			answer := app.AnswerPostParams{
				Content: "test content",
			}
			bytes, _ := json.Marshal(answer)
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			headers["token"] = getTestTok(&map[string]interface{}{"role": RoleAudience, "eid": 1, "qid": 1})
			req := performRequest(r, "POST", "/v1/answers/1/1", &bytes, headers)
			assert.Equal(t, http.StatusOK, req.Code)
		})
	})
}
