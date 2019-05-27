package controllers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/anon-solicitor/app"
	_ "github.com/mike-webster/anon-solicitor/app"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
)

var cfg = env.Config()

func getTestTok() string {
	id, _ := uuid.NewV4()
	return tokens.GetJWT(cfg.Secret, id.String())
}

func TestGetFeedbackV1(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		//TODO
	})

	t.Run("TokenError", func(t *testing.T) {
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		req := performRequest(r, "GET", "/v1/events/1/feedback/bad_token", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})

	t.Run("IdLessThanOne", func(t *testing.T) {
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/-1/feedback?token=%v", "cd625b90-1dde-4e76-a3d6-9ce8693ac6e1"), nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})

	t.Run("EventNotFound", func(t *testing.T) {
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/1000000/feedback?token=%v", "cd625b90-1dde-4e76-a3d6-9ce8693ac6e1"), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("GetFeedbackByTokError", func(t *testing.T) {
		opts := app.TestServiceOptions{
			ForceGetFeedbackByTokError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		tok := getTestTok()
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/1/feedback/%v?token=%v", tok, tok), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("FeedbackNotFound", func(t *testing.T) {
		opts := app.TestServiceOptions{
			ForceGetFeedbackByTokNotFound: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		tok := getTestTok()
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/1/feedback/%v?token=%v", tok, tok), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("Success", func(t *testing.T) {
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		tok := getTestTok()
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/1/feedback/%v?token=%v", tok, tok), nil, headers)
		assert.Equal(t, http.StatusOK, req.Code)
	})
}

func TestPostAbsentFeedbackV1(t *testing.T) {
	t.Run("DependenciesError", func(t *testing.T) {
		// TODO
	})

	t.Run("TokenError", func(t *testing.T) {
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, false)
		req := performRequest(r, "POST", "/v1/events/1/feedback/bad_token/absent", nil, headers)
		assert.Equal(t, http.StatusUnauthorized, req.Code)
	})

	t.Run("GetFeedbackByTokError", func(t *testing.T) {
		id, _ := uuid.NewV4()
		tok := tokens.GetJWT(cfg.Secret, id.String())
		opts := app.TestServiceOptions{
			ForceGetFeedbackByTokError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", fmt.Sprintf("/v1/events/1/feedback/%v/absent", tok), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code)
	})

	t.Run("MarkFeedbackAbsentError", func(t *testing.T) {
		tok := getTestTok()
		opts := app.TestServiceOptions{
			ForceMarkFeedbackAbsentError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", fmt.Sprintf("/v1/events/1/feedback/%v/absent", tok), nil, headers)
		assert.Equal(t, http.StatusInternalServerError, req.Code)
	})

	t.Run("Success", func(t *testing.T) {
		tok := getTestTok()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", fmt.Sprintf("/v1/events/1/feedback/%v/absent", tok), nil, headers)
		assert.Equal(t, http.StatusOK, req.Code)
	})
}

func TestPostFeedbackV1(t *testing.T) {
	// TODO
}
