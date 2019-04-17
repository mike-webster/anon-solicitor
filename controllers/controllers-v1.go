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
	r.PUT("/events/:id", putEventsV1)
	r.POST("/events/:id/feedback", postFeedbackV1)
	r.GET("/config", getConfigV1)
	r.PUT("/config", putConfigV1)

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
	es := c.Value(eventServiceKey.String())
	if es == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}
	c.HTML(200, "error.html", gin.H{"msg": "found service"})
	return
	// db, ok := contextDB.(DBWrapper)
	// if !ok {
	// 	c.HTML(500, "error.html", gin.H{"msg": "db conversion error"})
	// 	return
	// }

	// events := []Event{}
	// err := db.Get().Select(&events, "SELECT * FROM events")
	// if err != nil {
	// 	c.HTML(500, "error.html", gin.H{"msg": "db query error"})
	// 	return
	// }

	//c.HTML(http.StatusOK, "events.html", gin.H{"events": events})
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
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

func putEventsV1(c *gin.Context) {
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
	//contextDB := c.Value("db")
	// if contextDB == nil {
	// 	c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
	// 	return
	// }

	// db, ok := contextDB.(DBWrapper)
	// if !ok {
	// 	c.HTML(500, "error.html", gin.H{"msg": "db conversion error"})
	// 	return
	// }

	// eventUpdates := EventPutParams{}
	// err := c.Bind(&eventUpdates)
	// if err != nil {
	// 	log.Printf("Error binding object: %v", err)
	// 	c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": err})

	// 	return
	// }

	// existing := Event{}
	// db.Get().Where("id = ?", eventUpdates.ID).Find(&existing)

	// // TODO: Verify that the event selected from the database
	// //       belongs to the requesting user.  Maybe (after auth)
	// //       pull the user id from the token and add it to the
	// //       where clause?

	// // This isn't working because I'm not creating events correctly yet
	// // if existing.OrganizingUser.ID < 1 {
	// // 	c.HTML(http.StatusUnauthorized, "error.html", gin.H{"msg": err, "explain": "didn't find an event"})

	// // 	return
	// // }

	// if len(eventUpdates.Description) > 0 {
	// 	existing.Description = eventUpdates.Description
	// }

	// if len(eventUpdates.Title) > 0 {
	// 	existing.Title = eventUpdates.Title
	// }

	// if eventUpdates.OrganizerID > 0 {
	// 	newUser := User{}
	// 	db.Get().Where("id = ?", eventUpdates.OrganizerID).Find(&newUser)

	// 	if newUser.UserID < 1 {
	// 		c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": "new user wasn't found"})

	// 		return
	// 	}

	// 	existing.User = newUser
	// }

	// min := time.Time{}

	// if eventUpdates.Time != min {
	// 	existing.Time = eventUpdates.Time
	// }

	// db.Get().Save(&existing)

	// c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": existing})
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
