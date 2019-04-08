package main

import (
	"log"
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"
)

func getHomeV1(c *gin.Context) {
	c.HTML(http.StatusOK, "master.html", gin.H{"title": "Anon Solicitor"})
}

func getEventsV1(c *gin.Context) {
	contextDB := c.Value("db")
	if contextDB == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

	db, ok := contextDB.(DBWrapper)
	if !ok {
		c.HTML(500, "error.html", gin.H{"msg": "db conversion error"})
		return
	}

	events := []Event{}
	db.Get().First(&events)

	c.HTML(http.StatusOK, "events.html", gin.H{"events": events})
}

func postEventsV1(c *gin.Context) {
	contextDB := c.Value("db")
	if contextDB == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

	db, ok := contextDB.(DBWrapper)
	if !ok {
		c.HTML(500, "error.html", gin.H{"msg": "db conversion error"})
		return
	}

	newEvent := Event{}
	err := c.Bind(&newEvent)
	if err != nil {
		log.Printf("Error binding object: %v", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": err})

		return
	}

	db.Get().Create(&newEvent)

	c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": newEvent})
}

func putEventsV1(c *gin.Context) {
	contextDB := c.Value("db")
	if contextDB == nil {
		c.HTML(500, "error.html", gin.H{"msg": "missing db in context"})
		return
	}

	db, ok := contextDB.(DBWrapper)
	if !ok {
		c.HTML(500, "error.html", gin.H{"msg": "db conversion error"})
		return
	}

	eventUpdates := EventPutParams{}
	err := c.Bind(&eventUpdates)
	if err != nil {
		log.Printf("Error binding object: %v", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": err})

		return
	}

	existing := Event{}
	db.Get().Where("id = ?", eventUpdates.ID).Find(&existing)

	// TODO: Verify that the event selected from the database
	//       belongs to the requesting user.  Maybe (after auth)
	//       pull the user id from the token and add it to the
	//       where clause?

	// This isn't working because I'm not creating events correctly yet
	// if existing.OrganizingUser.ID < 1 {
	// 	c.HTML(http.StatusUnauthorized, "error.html", gin.H{"msg": err, "explain": "didn't find an event"})

	// 	return
	// }

	if len(eventUpdates.Description) > 0 {
		existing.Description = eventUpdates.Description
	}

	if len(eventUpdates.Name) > 0 {
		existing.Name = eventUpdates.Name
	}

	if eventUpdates.OrganizerID > 0 {
		newUser := User{}
		db.Get().Where("id = ?", eventUpdates.OrganizerID).Find(&newUser)

		if newUser.ID < 1 {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"msg": "new user wasn't found"})

			return
		}

		existing.OrganizingUser = newUser
	}

	min := time.Time{}

	if eventUpdates.Time != min {
		existing.Time = eventUpdates.Time
	}

	db.Get().Save(&existing)

	c.HTML(http.StatusOK, "event.html", gin.H{"title": "Anon Solicitor | Event", "event": existing})
}

func postFeedbackV1(c *gin.Context) {
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}

func getConfigV1(c *gin.Context) {
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}

func putConfigV1(c *gin.Context) {
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
}
