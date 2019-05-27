package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/mike-webster/anon-solicitor/app"
	domain "github.com/mike-webster/anon-solicitor/app"
	"github.com/mike-webster/anon-solicitor/data"
	"github.com/mike-webster/anon-solicitor/email"

	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/anon-solicitor/env"
)

var (
	testKey domain.ContextKey = "test"
)

const (
	controllerErrorKey      = "controllerError"
	controllerRespStatusKey = "responseStatus"
	ErrNoToken              = "err_missing_token"
	ErrBadToken             = "err_invalid_token"
	ErrUpdatingRecord       = "err_record_update"
	ErrNotImplemented       = "err_not_implemented"
)

// GetRouter will attempt to run the gin router
func GetRouter(ctx context.Context) *gin.Engine {
	cfg := env.Config()

	db, err := sqlx.Open("mysql", cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	createTables, ok := ctx.Value("DropTables").(bool)
	if !ok {
		createTables = false
	}
	if createTables {
		err = data.CreateTables(ctx, db)
		if err != nil {
			panic(err)
		}
	}

	return setupRouter(ctx, db)
}

func setupRouter(ctx context.Context, db *sqlx.DB) *gin.Engine {
	r := gin.Default()
	if env.Target() != "test" {
		r.LoadHTMLGlob("templates/*")
	}

	r.Use(setDependencies(ctx, db))
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
		v1Feedback.POST("/events/:id/feedback/:token/absent", postAbsentFeedbackV1)
	}

v1Question := r.Group("/v1"){
	v1Question.POST("/questions/:eventid/:token", postQuestionV1)
}

	// TODO: Catch all 404s
	// r.NoRoute(func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"test": "test"})
	// })

	return r
}

func getDependencies(ctx *gin.Context) (domain.EventService, domain.FeedbackService, domain.DeliveryService, error) {
	errs := ""

	es, err := getEventService(ctx, eventServiceKey.String())
	if err != nil {
		errs += err.Error() + ";"
	}

	fs, err := getFeedbackService(ctx, feedbackServiceKey.String())
	if err != nil {
		errs += err.Error() + ";"
	}

	em, err := getEmailService(ctx, domain.EmailServiceKey.String())

	if len(errs) > 1 {
		return nil, nil, nil, errors.New(errs)
	}

	return es, fs, em, nil
}

func getHomeV1(c *gin.Context) {
	c.HTML(http.StatusOK, "master.html", gin.H{"title": "Anon Solicitor"})
}

func setError(c *gin.Context, err error, desc string) {
	c.Error(gin.Error{
		Err:  err,
		Meta: desc,
	})

	c.Set(controllerErrorKey, true)
}

// getEmailService retrieves the expected EmailService with the give key from the gin context
func getEmailService(ctx *gin.Context, key interface{}) (domain.DeliveryService, error) {
	if ctx == nil {
		return nil, errors.New("provide a gin context in order to retrieve a value")
	}

	utEs := ctx.Value(key)
	if utEs == nil {
		return nil, errors.New("couldnt find key for Email Service in context")
	}

	tes, ok := utEs.(*app.TestDeliveryService)
	if ok {
		log.Print("warning: using test delivery service")
		return tes, nil
	}

	es, ok := utEs.(email.DeliveryService)
	if !ok {
		return nil, fmt.Errorf("couldnt parse Email Service from context, found: %v", reflect.TypeOf(utEs))
	}

	return &es, nil
}
