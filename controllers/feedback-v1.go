package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	gin "github.com/gin-gonic/gin"
	"github.com/mike-webster/anon-solicitor/app"
	domain "github.com/mike-webster/anon-solicitor/app"
	"github.com/mike-webster/anon-solicitor/data"
)

var feedbackServiceKey domain.ContextKey = "FeedbackService"

func postAbsentFeedbackV1(c *gin.Context) {
	_, fs, _, err := getDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, ErrRetrievingDependencies)

		return
	}

	tok, err := domain.String(c, "tok")
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusUnauthorized)
		setError(c, err, ErrBadToken)

		return
	} else if len(tok) < 1 {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusUnauthorized)
		setError(c, err, ErrNoToken)

		return
	}

	fb, err := fs.GetFeedbackByTok(tok)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, err, ErrRecordNotFound)

		return
	}

	err = fs.MarkFeedbackAbsent(fb)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, ErrUpdatingRecord)

		return
	}

	// TODO: update this to show a thanks for letting us know message
	c.HTML(http.StatusOK,
		"feedback.html",
		gin.H{"feedback": *fb})
}

func getFeedbackV1(c *gin.Context) {
	es, fs, _, err := getDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, ErrRetrievingDependencies)

		return
	}

	tok, err := domain.String(c, "tok")
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusUnauthorized)
		setError(c, err, ErrBadToken)

		return
	}

	eventid, _ := strconv.Atoi(c.Param("id"))
	if eventid < 1 {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusBadRequest)
		setError(c, errors.New("invalid event id"), ErrValidation)

		return
	}

	event := es.GetEvent(int64(eventid))

	if event == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, errors.New("event not found"), ErrRecordNotFound)

		return
	}

	fb, err := fs.GetFeedbackByTok(tok)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, err, ErrRecordNotFound)

		return
	}

	if fb == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, errors.New("event not found"), ErrRecordNotFound)

		return
	}

	// TODO: I need to add the questions to the feedback model
	c.HTML(http.StatusOK, "feedback.html", gin.H{"feedback": []domain.Feedback{*fb}})
}

func postFeedbackV1(c *gin.Context) {
	// TODO: Implement
	c.Set(controllerErrorKey, true)
	c.Set(controllerRespStatusKey, http.StatusNotImplemented)
	setError(c, errors.New("...coming soon..."), ErrNotImplemented)
}

// Feedback retrieves the expected EventService with the give key from the gin context
func getFeedbackService(ctx *gin.Context, key interface{}) (domain.FeedbackService, error) {
	if ctx == nil {
		return nil, errors.New("provide a gin context in order to retrieve a value")
	}

	tfs, ok := ctx.Value(key).(*app.TestFeedbackService)
	if ok {
		log.Print("warning: using test feedback service")
		return tfs, nil
	}

	fs, ok := ctx.Value(key).(data.FeedbackService)
	if !ok {
		return nil, fmt.Errorf("couldnt parse Feedback Service from context; found %v", reflect.TypeOf(ctx.Value(key)))
	}

	return fs, nil
}
