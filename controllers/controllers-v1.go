package controllers

import (
	"context"
	"log"
	"net/http"

	gin "github.com/gin-gonic/gin"
	anon "github.com/mike-webster/anon-solicitor"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

var eventServiceKey ContextKey = "EventService"
var testKey ContextKey = "test"

func StartServer(ctx context.Context, es anon.EventService) {
	r := setupRouter(ctx, es)
	r.Run("0.0.0.0:3001")
}

func setupRouter(ctx context.Context, es anon.EventService) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Use(setDependencies(es))

	r.GET("/", getHomeV1)
	r.GET("/events", getEventsV1)
	r.POST("/events", postEventsV1)
	r.POST("/events/:id/feedback", postFeedbackV1)

	// TODO: Catch all 404s
	// r.NoRoute(func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"test": "test"})
	// })

	return r
}

func setDependencies(es anon.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(eventServiceKey.String(), es)
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

	c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": posted})

	// if err != nil {
	// 	c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": err})
	// }

	// id, _ := res.LastInsertId()
	// newEvent.ID = id

	//c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": newEvent})
	//c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
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
