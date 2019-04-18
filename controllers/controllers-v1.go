package controllers

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	anon "github.com/mike-webster/anon-solicitor"
	"github.com/mike-webster/anon-solicitor/env"
	gomail "gopkg.in/gomail.v2"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

var eventServiceKey ContextKey = "EventService"
var feedbackServiceKey ContextKey = "FeedbackService"
var testKey ContextKey = "test"

func StartServer(ctx context.Context, es anon.EventService, fs anon.FeedbackService) {
	cfg := env.Config()
	r := setupRouter(ctx, es, fs)
	r.Run(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port))
}

func setupRouter(ctx context.Context, es anon.EventService, fs anon.FeedbackService) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Use(setDependencies(es, fs))

	r.GET("/", getHomeV1)
	r.GET("/events", getEventsV1)
	r.GET("/events/:id", getEventV1)
	r.POST("/events", postEventsV1)
	r.POST("/events/:id/feedback", postFeedbackV1)

	// TODO: Catch all 404s
	// r.NoRoute(func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"test": "test"})
	// })

	return r
}

func setDependencies(es anon.EventService, fs anon.FeedbackService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(eventServiceKey.String(), es)
		c.Set(feedbackServiceKey.String(), fs)
		c.Next()
	}
}

func getHomeV1(c *gin.Context) {
	c.HTML(http.StatusOK, "master.html", gin.H{"title": "Anon Solicitor"})
}

func getEventsV1(c *gin.Context) {
	untypedES := c.Value(eventServiceKey.String())
	if untypedES == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

	es, ok := untypedES.(anon.EventService)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": "couldnt cast db"})

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
	untypedES := c.Value(eventServiceKey.String())
	if untypedES == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

	es, ok := untypedES.(anon.EventService)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": "couldnt cast db"})

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
	untypedES := c.Value(eventServiceKey.String())
	if untypedES == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

	es, ok := untypedES.(anon.EventService)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": "couldnt cast db"})

		return
	}

	untypedFS := c.Value(feedbackServiceKey.String())
	if untypedFS == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

	fs, ok := untypedFS.(anon.FeedbackService)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": "couldnt cast db"})

		return
	}

	postEvent := anon.EventPostParams{}
	err := c.Bind(&postEvent)
	if err != nil {
		log.Printf("Error binding object: %v", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": err})

		return
	}

	log.Printf("posted event: %v", postEvent)

	newEvent := anon.Event{
		Title:       postEvent.Title,
		Description: postEvent.Description,
		Time:        postEvent.Time,
		UserID:      1,
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

	// send email to each audience member
	for k, v := range emails {
		err = sendEmail(v, k, posted.Title, posted.ID)
		if err != nil {
			log.Printf("Error: %v", err)

			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"msg": fmt.Sprintf("Error: %v", err)})
			return
		}
	}

	// probably redirect to the actual event page - which should show any submitted feedback
	// as well as how many total audiencemembers there were for the event

	c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": posted})
}

func sendEmail(email string, tok string, eventName string, eventID int64) error {
	cfg := env.Config()
	client := gomail.NewPlainDialer(cfg.SMTPHost, cfg.SMTPPort, "", "")
	client.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	message := gomail.NewMessage()
	message.SetHeader("From", "anno-solicitor@wyzant.com")
	message.SetHeader("To", email)

	fbPath := fmt.Sprintf("http://%v/events/%v/feedback/%v", cfg.Host, eventID, tok)
	body := fmt.Sprintf("<html><body><h3>Hey! We'd like to hear what you think!</h3><p>No worries - it's totally anonymous! Click <a href='%v'>here</a> to submit your feedback and see what everyone else thought!</p><p>Click <a href='%v'>here</a> to let us know that you didn't attend.</p><p>Thanks so much!</p></body></html>", fbPath, fbPath+"/absent")

	message.SetHeader("Title", fmt.Sprintf("You've been invited to give anonymous feedback about: %v", eventName))
	message.SetBody("text/html", body)

	if err := client.DialAndSend(message); err != nil {
		log.Printf("failed to send email. Error: %v", err)

		return err
	}

	return nil
}

func absentFeedbackV1(c *gin.Context) {

}

func getFeedbackV1(c *gin.Context) {

}

func postFeedbackV1(c *gin.Context) {
	// contextDB := c.Value("db")
	// if contextDB == nil {
	// 	c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
	// 	return
	// }

	// db, ok := contextDB.(DBWrapper)
	// if !ok {
	// 	c.HTML(500, "error.html", gin.H{"msg": "db conversion error"})
	// 	return
	// }

	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}

func getConfigV1(c *gin.Context) {
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}

func putConfigV1(c *gin.Context) {
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}
