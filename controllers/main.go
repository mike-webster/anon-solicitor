package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mike-webster/anon-solicitor/data"

	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	anon "github.com/mike-webster/anon-solicitor"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
	gomail "gopkg.in/gomail.v2"
)

var testKey anon.ContextKey = "test"

// StartServer will attempt to run the gin server
func StartServer(ctx context.Context) {
	cfg := env.Config()

	db, err := sqlx.Open("mysql", cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = data.CreateTables(ctx, db)
	if err != nil {
		panic(err)
	}

	r := setupRouter(ctx, db)

	r.Run(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port))
}

func setupRouter(ctx context.Context, db *sqlx.DB) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Use(setDependencies(ctx, db))

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

func setDependencies(ctx context.Context, db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		es := data.EventService{DB: db}
		fs := data.FeedbackService{DB: db}
		c.Set(eventServiceKey.String(), es)
		c.Set(feedbackServiceKey.String(), fs)
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

	es, err := getEventService(ctx, eventServiceKey.String())
	if err != nil {
		errs += err.Error() + ";"
	}

	fs, err := getFeedbackService(ctx, feedbackServiceKey.String())
	if err != nil {
		errs += err.Error() + ";"
	}

	if len(errs) > 1 {
		return nil, nil, errors.New(errs)
	}

	return es, fs, nil
}

// TODO: move this somewhere else?
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
