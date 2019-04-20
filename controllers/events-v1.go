package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	gin "github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	anon "github.com/mike-webster/anon-solicitor"
	"github.com/mike-webster/anon-solicitor/data"
	"github.com/mike-webster/anon-solicitor/env"
)

var eventServiceKey anon.ContextKey = "EventService"

func getEventsV1(c *gin.Context) {
	es, _, err := getDependencies(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	events, err := es.GetEvents()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": "db query error", "err": err})

		return
	}

	c.HTML(http.StatusOK, "events.html", gin.H{"events": events})
}

func getEventV1(c *gin.Context) {
	es, _, err := getDependencies(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

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

	c.HTML(http.StatusOK, "events.html", gin.H{"events": []anon.Event{*event}})
}

func postEventsV1(c *gin.Context) {
	es, fs, err := getDependencies(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError,
			"error.html",
			gin.H{"msg": err})

		return
	}

	postEvent := anon.EventPostParams{}
	err = c.Bind(&postEvent)
	if err != nil {
		log.Printf("Error binding object: %v", err)
		c.HTML(http.StatusBadRequest,
			"error.html",
			gin.H{"msg": err})

		return
	}

	log.Printf("posted event: %v", postEvent)

	newEvent := anon.Event{
		Title:       postEvent.Title,
		Description: postEvent.Description,
		Time:        postEvent.Time,
	}

	log.Printf("saving event: %v", newEvent)

	err = es.CreateEvent(&newEvent)
	if err != nil {
		log.Printf("error creating event: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": err})

		return
	}

	posted := es.GetEvent(newEvent.ID)
	if posted == nil {
		log.Printf("error getting created event: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": "couldnt retrieve saved event", "id": newEvent.ID})

		return
	}

	if env.Config().ShouldSendEmails {

		emails := map[string]string{}

		for _, email := range postEvent.Audience {
			// create feedback record for each audience member
			// - attach tok to each one
			tok, err := uuid.NewV4()
			if err != nil {
				c.HTML(http.StatusInternalServerError,
					"error.html",
					gin.H{"msg": "problem creating tokens", "id": newEvent.ID})

				return
			}

			emails[tok.String()] = email

			newFeedback := anon.Feedback{
				Tok:     tok.String(),
				EventID: posted.ID,
			}

			// clear the tok when the feedback is submitted
			err = fs.CreateFeedback(&newFeedback)
			if err != nil {
				c.HTML(http.StatusInternalServerError,
					"error.html",
					gin.H{
						"msg": fmt.Sprintf("problem creating tokens, Err: %v", err),
						"id":  newEvent.ID,
					})

				return
			}
		}

		// TODO: test this part
		// send email to each audience member
		for k, v := range emails {
			err = sendEmail(v, k, posted.Title, posted.ID)
			if err != nil {
				log.Printf("Error: %v", err)

				c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": fmt.Sprintf("Error: %v", err)})
				return
			}
		}
	}

	// probably redirect to the actual event page - which should show any submitted feedback
	// as well as how many total audiencemembers there were for the event

	c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": posted})
}

// getEventService retrieves the expected EventService with the give key from the gin context
func getEventService(ctx *gin.Context, key interface{}) (anon.EventService, error) {
	if ctx == nil {
		return nil, errors.New("provide a gin context in order to retrieve a value")
	}

	utEs := ctx.Value(key)
	if utEs == nil {
		return nil, errors.New("couldnt find key for Event Service in context")
	}

	es, ok := utEs.(data.EventService)
	if !ok {
		return nil, fmt.Errorf("couldnt parse Event Service from context, found: %v", reflect.TypeOf(utEs))
	}

	return &es, nil
}
