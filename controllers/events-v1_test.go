package controllers

import (
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	gin "github.com/gin-gonic/gin"
	"github.com/mike-webster/anon-solicitor/app"
)

// TODO: restructure this so we use one method for both
func setupTestRouter(deps *app.AnonDependencies) *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("../templates/*")
	r.Use(func(c *gin.Context) {
		c.Set("EventService", deps.Events)
		c.Set("FeedbackService", deps.Feedback)
		c.Next()
	})
	r.Use(setStatus())

	v1Events := r.Group("/v1")
	{
		v1Events.GET("/", getHomeV1)
		v1Events.GET("/events", getEventsV1)
		v1Events.GET("/events/:id", getEventV1)
		v1Events.POST("/events", postEventsV1)
	}

	r.Use(getToken())

	// TODO: isolate these into a group so I can use the getToken()
	//       middleware on only these routes.
	// TODO: make sure this doesn't cause any weirdness... i'm delcaring a
	//       second "/v1" on this router
	v1Feedback := r.Group("/v1")
	{
		v1Feedback.GET("/events/:id/feedback/:token", getFeedbackV1)
		v1Feedback.POST("/events/:id/feedback/:token", postFeedbackV1)
		v1Feedback.POST("/events/:id/feedback/:token/absent", postAbsentFeedbackV1)
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
		opts := app.TestServiceOptions{
			ForceGetEventsError: true,
		}
		deps := app.MockSearchDependencies(opts)
		r := setupTestRouter(deps)
		req := performRequest(r, "GET", "/v1/events")
		assert.Equal(t, http.StatusInternalServerError, req.Code)
	})

	t.Run("Success", func(t *testing.T) {
		deps := app.MockSearchDependencies(app.TestServiceOptions{})
		r := setupTestRouter(deps)
		req := performRequest(r, "GET", "/v1/events")
		assert.Equal(t, http.StatusOK, req.Code)
	})
}

func TestGetEventV1(t *testing.T) {

	r := gin.New()
	r.GET("/v1/event", getEventsV1)

	t.Run("TestNoDependencies", func(t *testing.T) {
		// TODO: figure out how to make this happen
	})

	t.Run("TestIDLessThan1Invalid", func(t *testing.T) {
	})

	t.Run("TestIDNotFound", func(t *testing.T) {
	})

	t.Run("TestSuccess", func(t *testing.T) {
	})
}

func TestPostEventV1(t *testing.T) {

	r := gin.New()
	r.GET("/v1/event", getEventsV1)

	t.Run("TestNoDependencies", func(t *testing.T) {
		// TODO: figure out how to make this happen
	})

	t.Run("TestValidation", func(t *testing.T) {
		t.Run("TitleNotProvided", func(t *testing.T) {

		})
		t.Run("TitleLongerThanTwoHundredCharacters", func(t *testing.T) {

		})
		t.Run("DescriptionNotProvided", func(t *testing.T) {

		})
		t.Run("DescriptionLongerThanFiveThousandCharacters", func(t *testing.T) {

		})
		t.Run("ScheduledTimeNotProvided", func(t *testing.T) {

		})
		t.Run("AudienceNotProvided", func(t *testing.T) {

		})
	})

	t.Run("TestCreationError", func(t *testing.T) {

	})

	t.Run("TestVerificationError", func(t *testing.T) {

	})

	t.Run("TestShouldSendEmails", func(t *testing.T) {
		t.Run("TestCreateFeedbackError", func(t *testing.T) {

		})
	})

	t.Run("TestSuccess", func(t *testing.T) {

	})
}
