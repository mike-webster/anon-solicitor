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
	"github.com/mike-webster/anon-solicitor/app"
	domain "github.com/mike-webster/anon-solicitor/app"
	"github.com/mike-webster/anon-solicitor/data"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
)

var (
	eventServiceKey domain.ContextKey = "EventService"
)

const (
	ErrRetrievingDependencies = "error_retrieving_dependencies"
	ErrRetrievingDomainObject = "error_retrieving_domain_object"
	ErrValidation             = "error_validation"
	ErrRecordNotFound         = "err_record_not_found"
	ErrRecordCreation         = "err_record_not_created"
	ErrTokenCreation          = "err_token_not_created"
	ErrEmail                  = "err_sending_email"
)

func getEventsV1(c *gin.Context) {
	es, _, _, err := getDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		setError(c, err, ErrRetrievingDependencies)

		return
	}

	events, err := es.GetEvents()
	if err != nil {
		c.Set(controllerErrorKey, true)
		setError(c, err, ErrRetrievingDomainObject)

		return
	}

	c.HTML(http.StatusOK, "events.html", gin.H{"events": events})
}

func getEventV1(c *gin.Context) {
	es, _, _, err := getDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		setError(c, err, ErrRetrievingDependencies)

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

	c.HTML(http.StatusOK, "events.html", gin.H{"events": []domain.Event{*event}})
}

func postEventsV1(c *gin.Context) {
	cfg := env.Config()
	es, fs, em, err := getDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		setError(c, err, ErrRetrievingDependencies)

		return
	}

	postEvent := domain.EventPostParams{}
	err = c.Bind(&postEvent)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusBadRequest)
		setError(c, err, ErrValidation)
		log.Printf("Error binding object: %v", err)

		return
	}

	log.Printf("posted event: %v", postEvent)

	newEvent := domain.Event{
		Title:       postEvent.Title,
		Description: postEvent.Description,
		Time:        postEvent.Time,
	}

	log.Printf("saving event: %v", newEvent)

	err = es.CreateEvent(&newEvent)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, ErrRecordCreation)
		log.Printf("error creating event: %v", err)

		return
	}

	posted := es.GetEvent(newEvent.ID)
	if posted == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		err := fmt.Errorf("couldnt find newly created event - id: %v", newEvent.ID)
		setError(c, err, ErrRecordNotFound)
		log.Printf("error getting created event: %v", err)

		return
	}

	if env.Config().ShouldSendEmails {

		emails := map[string]string{}

		for _, email := range postEvent.Audience {
			// create feedback record for each audience member
			// - attach tok to each one
			tok, err := uuid.NewV4()
			if err != nil {
				c.Set(controllerErrorKey, true)
				c.Set(controllerRespStatusKey, http.StatusInternalServerError)
				setError(c, err, ErrTokenCreation)

				return
			}

			emails[tok.String()] = email

			newFeedback := domain.Feedback{
				Tok:     tok.String(),
				EventID: posted.ID,
			}

			// clear the tok when the feedback is submitted
			err = fs.CreateFeedback(&newFeedback)
			if err != nil {
				c.Set(controllerErrorKey, true)
				c.Set(controllerRespStatusKey, http.StatusInternalServerError)
				setError(c, err, ErrRecordCreation)

				return
			}
		}

		// TODO: test this part
		// send email to each audience member
		for k, v := range emails {
			fbPath := fmt.Sprintf("http://%v/events/%v/feedback/%v",
				cfg.Host,
				posted.ID,
				tokens.GetJWT(cfg.Secret, k))

			err = em.SendFeedbackEmail(v, fbPath)
			if err != nil {
				c.Set(controllerErrorKey, true)
				c.Set(controllerRespStatusKey, http.StatusInternalServerError)
				setError(c, err, ErrEmail)
				log.Printf("Error: %v", err)

				return
			}
		}
	}

	// probably redirect to the actual event page - which should show any submitted feedback
	// as well as how many total audiencemembers there were for the event

	c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": posted})
}

// getEventService retrieves the expected EventService with the give key from the gin context
func getEventService(ctx *gin.Context, key interface{}) (domain.EventService, error) {
	if ctx == nil {
		return nil, errors.New("provide a gin context in order to retrieve a value")
	}

	utEs := ctx.Value(key)
	if utEs == nil {
		return nil, errors.New("couldnt find key for Event Service in context")
	}

	tes, ok := utEs.(*app.TestEventService)
	if ok {
		log.Print("warning: using test event service")
		return tes, nil
	}

	es, ok := utEs.(data.EventService)
	if !ok {
		return nil, fmt.Errorf("couldnt parse Event Service from context, found: %v", reflect.TypeOf(utEs))
	}

	return &es, nil
}
