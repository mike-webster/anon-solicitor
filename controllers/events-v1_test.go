package controllers

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	gin "github.com/gin-gonic/gin"
	"github.com/mike-webster/anon-solicitor/app"
	domain "github.com/mike-webster/anon-solicitor/app"
	"github.com/mike-webster/anon-solicitor/env"
)

func getTestHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}

// TODO: restructure this so we use one method for both
func setupTestRouter(deps *app.AnonDependencies, useAuth bool) *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("../templates/*")
	r.Use(func(c *gin.Context) {
		c.Set("EventService", deps.Events)
		c.Set("FeedbackService", deps.Feedback)
		c.Set("EmailService", deps.Delivery)
		c.Next()
	})

	r.GET("/testsetup", getTestSetup)
	r.POST("/testsetup", postTestSetup)

	r.Use(setStatus())

	v1Events := r.Group("/v1")
	{
		v1Events.GET("/", getHomeV1)
		v1Events.GET("/events", getEventsV1)
		v1Events.GET("/events/:id", getEventV1)
		v1Events.POST("/events", postEventsV1)
	}

	if useAuth {
		r.Use(getToken())
	}
	// TODO: isolate these into a group so I can use the getToken()
	//       middleware on only these routes.
	// TODO: make sure this doesn't cause any weirdness... i'm delcaring a
	//       second "/v1" on this router
	v1Feedback := r.Group("/v1")
	{
		v1Feedback.GET("/events/:id/feedback", getFeedbackV1)
		v1Feedback.POST("/events/:id/feedback/absent", postAbsentFeedbackV1)
	}

	v1Question := r.Group("/v1")
	{
		v1Question.POST("/questions/:eventid", postQuestionV1)
	}

	v1Answer := r.Group("/v1")
	{
		v1Answer.POST("/answers/:eventid/:questionid", postQuestionAnswerV1)
	}
	// TODO: Catch all 404s
	// r.NoRoute(func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"test": "test"})
	// })

	return r
}

func TestGetEventsV1(t *testing.T) {
	t.Run("TestNoDependencies", func(t *testing.T) {
		// TODO: figure out how to make this happen
	})

	t.Run("TestErrorRetrievingEvents", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceGetEventsError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", "/v1/events", nil, headers)
		assert.Equal(t, http.StatusInternalServerError, req.Code)
	})

	t.Run("Success", func(t *testing.T) {
		headers := getTestHeaders()
		deps := app.MockSearchDependencies(app.TestServiceOptions{})
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", "/v1/events", nil, headers)
		assert.Equal(t, http.StatusOK, req.Code)
	})
}

func TestGetEventV1(t *testing.T) {
	t.Run("TestNoDependencies", func(t *testing.T) {
		// TODO: figure out how to make this happen
	})

	t.Run("TestIDLessThan1Invalid", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", "/v1/events/0", nil, headers)
		assert.Equal(t, http.StatusBadRequest, req.Code, req.Body.String())
	})

	t.Run("TestIDNotFound", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{
			ForceGetEventError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", fmt.Sprintf("/v1/events/%v", 30000), nil, headers)
		assert.Equal(t, http.StatusNotFound, req.Code, req.Body.String())
	})

	t.Run("TestSuccess", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "GET", "/v1/events/1", nil, headers)
		assert.Equal(t, http.StatusOK, req.Code, req.Body.String())
	})
}

func getValidEventParams() domain.EventPostParams {
	return domain.EventPostParams{
		Title:       "Test Title",
		Description: "Test Description",
		Time:        time.Now(),
		Audience:    []string{"test@testemail.com"},
	}
}

func TestPostEventV1(t *testing.T) {
	t.Run("TestNoDependencies", func(t *testing.T) {
		// TODO: figure out how to make this happen
	})
	t.Run("TestValidation", func(t *testing.T) {
		t.Run("TitleNotProvided", func(t *testing.T) {
			headers := getTestHeaders()
			e := getValidEventParams()
			e.Title = ""
			b, _ := json.Marshal(e)
			opts := app.TestServiceOptions{}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body.String())

			t.Run("ExpectedError", func(t *testing.T) {
				assert.Equal(t, true, strings.Contains(req.Body.String(), "{\"Title\":\"Title is required\"}"), req.Body.String())
			})
		})
		t.Run("TitleLongerThanTwoHundredCharacters", func(t *testing.T) {
			headers := getTestHeaders()
			e := getValidEventParams()
			for i := 0; i < 201; i++ {
				e.Title += "a"
			}
			log.Println("\n\n\n~~~len: ", len(e.Description))
			b, _ := json.Marshal(e)
			opts := app.TestServiceOptions{}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body.String())

			t.Run("ExpectedError", func(t *testing.T) {
				assert.Equal(t, true, strings.Contains(req.Body.String(), "{\"Title\":\"Title cannot be longer than 200\"}"), req.Body.String())
			})
		})
		t.Run("DescriptionNotProvided", func(t *testing.T) {
			headers := getTestHeaders()
			e := getValidEventParams()
			e.Description = ""
			b, _ := json.Marshal(e)
			opts := app.TestServiceOptions{}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body.String())

			t.Run("ExpectedError", func(t *testing.T) {
				assert.Equal(t, true, strings.Contains(req.Body.String(), "{\"Description\":\"Description is required\"}"), req.Body.String())
			})
		})
		t.Run("DescriptionLongerThanFiveThousandCharacters", func(t *testing.T) {
			headers := getTestHeaders()
			e := getValidEventParams()
			e.Description = ""
			for i := 0; i < 5001; i++ {
				e.Description += "a"
			}
			b, _ := json.Marshal(e)
			opts := app.TestServiceOptions{}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body.String())

			t.Run("ExpectedError", func(t *testing.T) {
				assert.Equal(t, true, strings.Contains(req.Body.String(), "{\"Description\":\"Description cannot be longer than 5000\"}"), req.Body.String())
			})
		})
		t.Run("ScheduledTimeNotProvided", func(t *testing.T) {
			headers := getTestHeaders()
			e := getValidEventParams()
			e.Time = time.Time{}
			b, _ := json.Marshal(e)
			opts := app.TestServiceOptions{}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body.String())

			t.Run("ExpectedError", func(t *testing.T) {
				assert.Equal(t, true, strings.Contains(req.Body.String(), "{\"Time\":\"Time is required\"}"), req.Body.String())
			})
		})
		t.Run("AudienceNotProvided", func(t *testing.T) {
			headers := getTestHeaders()
			e := getValidEventParams()
			e.Audience = []string{}
			b, _ := json.Marshal(e)
			opts := app.TestServiceOptions{}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusBadRequest, req.Code, req.Body.String())

			t.Run("ExpectedError", func(t *testing.T) {
				assert.Equal(t, true, strings.Contains(req.Body.String(), "{\"Audience\":\"Audience must be longer than 1\"}"), req.Body.String())
			})
		})
	})
	t.Run("TestCreationError", func(t *testing.T) {
		headers := getTestHeaders()
		e := getValidEventParams()
		b, _ := json.Marshal(e)
		opts := app.TestServiceOptions{
			ForceCreateEventError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", "/v1/events", &b, headers)
		assert.Equal(t, http.StatusInternalServerError, req.Code, req.Body.String())

		t.Run("ExpectedError", func(t *testing.T) {
			assert.Equal(t, true, strings.Contains(req.Body.String(), "forced test error"), req.Body.String())
		})
	})
	t.Run("TestVerificationError", func(t *testing.T) {
		headers := getTestHeaders()
		e := getValidEventParams()
		b, _ := json.Marshal(e)
		opts := app.TestServiceOptions{
			ForceGetEventError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", "/v1/events", &b, headers)
		assert.Equal(t, http.StatusInternalServerError, req.Code, req.Body.String())

		t.Run("ExpectedError", func(t *testing.T) {
			assert.NotEqual(t, http.StatusOK, req.Code)
			assert.Equal(t, true, strings.Contains(req.Body.String(), "couldnt find newly created event - id:"), req.Body.String())
		})
	})
	t.Run("TestShouldSendEmails", func(t *testing.T) {
		t.Run("SkipsSendingIfConfigured", func(t *testing.T) {
			headers := getTestHeaders()
			os.Setenv("SEND_EMAILS", "false")
			t.Run("ConfiguredCorrectly", func(t *testing.T) {
				assert.Equal(t, "false", os.Getenv("SEND_EMAILS"))
			})

			b, _ := json.Marshal(getValidEventParams())
			deps := app.MockSearchDependencies(app.TestServiceOptions{})
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusOK, req.Code, req.Body.String())
			t.Run("EmailsWerentSent", func(t *testing.T) {
				ds, _ := deps.Delivery.(*app.TestDeliveryService)
				assert.Equal(t, 0, ds.GetFeedbackEmailCount())
			})
		})
		t.Run("TestCreateFeedbackError", func(t *testing.T) {
			headers := getTestHeaders()
			os.Setenv("SEND_EMAILS", "true")
			t.Run("ConfiguredCorrectly", func(t *testing.T) {
				assert.Equal(t, "true", os.Getenv("SEND_EMAILS"))
			})

			b, _ := json.Marshal(getValidEventParams())
			opts := app.TestServiceOptions{
				ForceCreateFeedbackError: true,
			}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, true)
			req := performRequest(r, "POST", "/v1/events", &b, headers)
			assert.Equal(t, http.StatusOK, req.Code, req.Body.String())
			t.Run("EmailsWerentSent", func(t *testing.T) {
				ds, _ := deps.Delivery.(*app.TestDeliveryService)
				assert.Equal(t, 0, ds.GetFeedbackEmailCount())
			})
		})
	})

	t.Run("TestSuccess", func(t *testing.T) {
		headers := getTestHeaders()
		cfg := env.Config()
		cfg.ShouldSendEmails = true

		e := getValidEventParams()
		b, _ := json.Marshal(e)
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps, true)
		req := performRequest(r, "POST", "/v1/events", &b, headers)
		assert.Equal(t, http.StatusOK, req.Code, req.Body.String())
		t.Run("EmailsWereSent", func(t *testing.T) {
			ds, _ := deps.Delivery.(*app.TestDeliveryService)
			assert.Equal(t, 1, ds.GetFeedbackEmailCount())
		})
	})
}
