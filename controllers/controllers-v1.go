package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mike-webster/anon-solicitor/data"

	gin "github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	anon "github.com/mike-webster/anon-solicitor"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
	gomail "gopkg.in/gomail.v2"
)

var eventServiceKey anon.ContextKey = "EventService"
var feedbackServiceKey anon.ContextKey = "FeedbackService"
var testKey anon.ContextKey = "test"

// StartServer will attempt to run the gin server
func StartServer(ctx context.Context) {
	cfg := env.Config()

	db, err := sqlx.Open("mysql", cfg.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := setupRouter(ctx)

	err = data.CreateTables(ctx, db)
	if err != nil {
		panic(err)
	}

	r.Run(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port))
}

func setupRouter(ctx context.Context) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Use(setDependencies(ctx))

	r.GET("/", getHomeV1)
	r.GET("/events", getEventsV1)
	r.GET("/events/:id", getEventV1)
	r.POST("/events", postEventsV1)

	r.Use(getToken())

	r.GET("/events/:id/feedback/:token", getFeedbackV1)
	r.POST("/events/:id/feedback/:token", postFeedbackV1)
	r.POST("/events/:id/feedback/:token/absent", absentFeedbackV1)

	// TODO: Catch all 404s
	// r.NoRoute(func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"test": "test"})
	// })

	return r
}

func setDependencies(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := anon.DB(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		es := data.EventService{DB: db}
		fs := data.FeedbackService{DB: db}
		c.Set(eventServiceKey.String(), &es)
		c.Set(feedbackServiceKey.String(), &fs)
		c.Next()
	}
}

func getToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Note: this is probably unncessary if the token is going to be a url param...
		//       I just wanted to do it. :)
		// TODO: test this
		cfg := env.Config()
		token := c.Param("token")
		if len(token) < 1 {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		tok, err := tokens.CheckToken(token, cfg.Secret)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		if len(tok) < 1 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("couldn't find token"))

			return
		}

		c.Set("tok", tok)
		c.Next()
	}
}

func getDependencies(ctx *gin.Context) (anon.EventService, anon.FeedbackService, error) {

	errs := ""

	es, err := anon.GetEventService(ctx, eventServiceKey.String())
	if err != nil {
		errs += err.Error() + ";"
	}

	fs, err := anon.GetFeedbackService(ctx, feedbackServiceKey.String())
	if err != nil {
		errs += err.Error() + ";"
	}

	if len(errs) > 1 {
		return nil, nil, errors.New(errs)
	}

	return *es, *fs, nil
}

func sendEmail(email string, tok string, eventName string, eventID int64) error {
	cfg := env.Config()
	client := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	message := gomail.NewMessage()
	message.SetHeader("From", fmt.Sprintf("Anon Solicitor <%v>", cfg.SMTPUser))
	message.SetHeader("To", email)
	jwt := tokens.GetJWT(cfg.Secret, tok)
	fbPath := fmt.Sprintf("http://%v/events/%v/feedback/%v", cfg.Host, 1, jwt)
	body := fmt.Sprintf("<html><body><h3>Hey! We'd like to hear what you think!</h3><p>No worries - it's totally anonymous! Click <a href='%v'>here</a> to submit your feedback and see what everyone else thought!</p><p>Click <a href='%v'>here</a> to let us know that you didn't attend.</p><p>Thanks so much!</p></body></html>", fbPath, fbPath+"/absent")

	message.SetHeader("Title", fmt.Sprintf("You've been invited to give anonymous feedback about: %v", "test event"))
	message.SetBody("text/html", body)

	if err := client.DialAndSend(message); err != nil {
		log.Printf("failed to send email. Error: %v", err)

		return err
	}

	return nil
}

func getHomeV1(c *gin.Context) {
	c.HTML(http.StatusOK, "master.html", gin.H{"title": "Anon Solicitor"})
}

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

	err = fs.MarkFeebackAbsent(fb)
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
	// TODO: Implement
}

func postFeedbackV1(c *gin.Context) {
	// TODO: Implement
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}
