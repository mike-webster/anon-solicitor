package controllers

import (
	"context"
	"net/http"

	gin "github.com/gin-gonic/gin"
)

func StartServer(ctx context.Context) {
	r := setupRouter(ctx)
	r.Run("0.0.0.0:3001")
}

func setupRouter(ctx context.Context) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	//r.Use(mwAttachDB)

	r.GET("/", getHomeV1)
	r.GET("/events", getEventsV1)
	r.POST("/events", postEventsV1)
	r.PUT("/events/:id", putEventsV1)
	r.POST("/events/:id/feedback", postFeedbackV1)
	r.GET("/config", getConfigV1)
	r.PUT("/config", putConfigV1)

	return r
}

func getHomeV1(c *gin.Context) {
	c.HTML(http.StatusOK, "master.html", gin.H{"title": "Anon Solicitor"})
}

func getEventsV1(c *gin.Context) {
	contextDB := c.Value("db")
	if contextDB == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

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

	// newEvent := Event{}
	// err := c.Bind(&newEvent)
	// if err != nil {
	// 	log.Printf("Error binding object: %v", err)
	// 	c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": err})

	// 	return
	// }

	// createdAt := time.Now().UTC()
	// res, err := db.Get().Exec("INSERT INTO events (id, title, description, time, created_at, modified_at) VALUES (?,?,?,?,?,?)",
	// 	newEvent.ID,
	// 	newEvent.Title,
	// 	newEvent.Description,
	// 	newEvent.Time,
	// 	createdAt,
	// 	createdAt,
	// )

	// if err != nil {
	// 	c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": err})
	// }

	// id, _ := res.LastInsertId()
	// newEvent.ID = id

	//c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": newEvent})
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
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
