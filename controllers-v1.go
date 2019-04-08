package main

import (
	"log"
	"net/http"

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
	c.HTML(http.StatusNotImplemented, "error.html", gin.H{"msg": "...coming soon..."})
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
