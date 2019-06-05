package controllers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/anon-solicitor/app"
	_ "github.com/mike-webster/anon-solicitor/app"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
)

var cfg = env.Config()

func getTestTok(inPayload *map[string]interface{}) string {
	id, _ := uuid.NewV4()
	foundTok := false
	payload := map[string]interface{}{}
	if inPayload != nil {
		for k, v := range *inPayload {
			payload[k] = v
			if k == "tok" {
				foundTok = true
			}
		}
	}
	if !foundTok {
		payload["tok"] = id.String()
	}
	payload["exp"] = time.Now().UTC().Add(30 * time.Minute).Unix()
	payload["iss"] = "anon-test"
	return tokens.GetJWT(cfg.Secret, payload)
}

func TestGetFeedbackV1(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		//TODO
	})

	t.Run("TokenError", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		req := performRequest(r, "GET", "/v1/events/1/feedback", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})

	t.Run("IdLessThanOne", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", "/v1/events/-1/feedback", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})

	t.Run("EventNotFound", func(t *testing.T) {
		headers := getTestHeaders()
		headers["token"] = getTestTok(nil)
		opts := app.TestServiceOptions{
			ForceGetEventError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", "/v1/events/1000000/feedback", nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("GetFeedbackByTokError", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceGetFeedbackByTokError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		tok := getTestTok(nil)
		headers["token"] = tok
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/1/feedback/%v?token=%v", tok, tok), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("FeedbackNotFound", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceGetFeedbackByTokNotFound: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		tok := getTestTok(nil)
		headers["token"] = tok
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/1/feedback/%v?token=%v", tok, tok), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("Success", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		tok := getTestTok(nil)
		headers["token"] = tok
		req := performRequest(r, "GET", "/v1/events/1/feedback", nil, headers)
		assert.Equal(t, http.StatusOK, req.Code)
	})
}

func TestPostAbsentFeedbackV1(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		// TODO
	})

	t.Run("TokenError", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		req := performRequest(r, "POST", "/v1/events/1/feedback/absent", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})

	t.Run("GetFeedbackByTokError", func(t *testing.T) {
		headers := getTestHeaders()
		tok := getTestTok(nil)
		headers["token"] = tok
		opts := app.TestServiceOptions{
			ForceGetFeedbackByTokError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", fmt.Sprintf("/v1/events/1/feedback/%v/absent", tok), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("MarkFeedbackAbsentError", func(t *testing.T) {
		headers := getTestHeaders()
		tok := getTestTok(nil)
		headers["token"] = tok
		opts := app.TestServiceOptions{
			ForceMarkFeedbackAbsentError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", "/v1/events/1/feedback/absent", nil, headers)
		assert.Equal(t, http.StatusInternalServerError, req.Code)
	})

	t.Run("Success", func(t *testing.T) {
		headers := getTestHeaders()
		tok := getTestTok(nil)
		headers["token"] = tok
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", "/v1/events/1/feedback/absent", nil, headers)
		assert.Equal(t, http.StatusOK, req.Code)
	})
}
