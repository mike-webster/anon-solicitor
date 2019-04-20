package controllers

import (
	"errors"
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"
	anon "github.com/mike-webster/anon-solicitor"
	"github.com/mike-webster/anon-solicitor/data"
)

var feedbackServiceKey anon.ContextKey = "FeedbackService"

func absentFeedbackV1(c *gin.Context) {
	_, fs, err := getDependencies(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError,
			"error.html",
			gin.H{"msg": err})

		return
	}

	tok, err := anon.String(c, "tok")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	} else if len(tok) < 1 {
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	fb, err := fs.GetFeedbackByTok(tok)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	err = fs.MarkFeedbackAbsent(fb)
	if err != nil {
		c.HTML(http.StatusInternalServerError,
			"error.html",
			gin.H{"msg": err})

		return
	}

	// TODO: update this to show a thanks for letting us know message
	c.HTML(http.StatusOK,
		"feedback.html",
		gin.H{"feedback": *fb})
}

func getFeedbackV1(c *gin.Context) {
	es, fs, err := getDependencies(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	tok, err := anon.String(c, "tok")
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)

		return
	}

	eventid, _ := strconv.Atoi(c.Param("id"))
	if eventid < 1 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": "invalid event id"})

		return
	}

	event := es.GetEvent(int64(eventid))

	if event == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"msg": "event not found"})

		return
	}

	fb, err := fs.GetFeedbackByTok(tok)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	if fb == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"msg": "feedback record not found"})

		return
	}

	// TODO: I need to add the questions to the feedback model
	c.HTML(http.StatusOK, "feedback.html", gin.H{"feedback": []anon.Feedback{*fb}})
}

func postFeedbackV1(c *gin.Context) {
	// TODO: Implement
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}

// Feedback retrieves the expected EventService with the give key from the gin context
func getFeedbackService(ctx *gin.Context, key interface{}) (anon.FeedbackService, error) {

	if ctx == nil {
		return nil, errors.New("provide a gin context in order to retrieve a value")
	}

	fs, ok := ctx.Value(key).(data.FeedbackService)
	if !ok {
		return nil, errors.New("couldnt parse Feedback Service from context")
	}

	return fs, nil
}
