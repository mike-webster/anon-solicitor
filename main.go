package main

import (
	gin "github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	r.Run("0.0.0.0:3001") // listen and serve on 0.0.0.0:8080
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", getHomeV1)
	r.GET("/events", getEventsV1)
	r.POST("/events", postEventsV1)
	r.PUT("/events/:id", putEventsV1)
	r.POST("/events/:id/feedback", postFeedbackV1)
	r.GET("/config", getConfigV1)
	r.PUT("/config", putConfigV1)

	return r
}
