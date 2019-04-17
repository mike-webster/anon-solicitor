package main

import (
	"context"
	"log"
	"time"

	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mike-webster/anon-solicitor/data"
)

func main() {
	ctx := context.Background()
	log.Print("Sleeping to allow db setup...")
	time.Sleep(10 * time.Second)

	wrap := DBWrapper{}
	wrap.Get()

	ctx = context.WithValue(ctx, "db", wrap)
	err := data.CreateTables(ctx)
	if err != nil {
		panic(err)
	}

	r := setupRouter()
	r.Run("0.0.0.0:3001")
}

func mwAttachDB(c *gin.Context) {
	wrap := DBWrapper{}
	wrap.Get()
	c.Set("db", wrap)

	defer func() {
		if r := recover(); r != nil {
			log.Println("Middleware caught a panic", r)
		}
	}()

	c.Next()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Use(mwAttachDB)

	r.GET("/", getHomeV1)
	r.GET("/events", getEventsV1)
	r.POST("/events", postEventsV1)
	r.PUT("/events/:id", putEventsV1)
	r.POST("/events/:id/feedback", postFeedbackV1)
	r.GET("/config", getConfigV1)
	r.PUT("/config", putConfigV1)

	return r
}
