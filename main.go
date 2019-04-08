package main

import (
	"log"
	"time"

	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	log.Print("Sleeping to allow db setup...")
	time.Sleep(2 * time.Second)

	r := setupRouter()
	r.Run("0.0.0.0:3001")
}

func mwAttachDB(c *gin.Context) {
	c.Set("db", DBWrapper{db()})

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

func db() *gorm.DB {
	db, err := gorm.Open("mysql", "root@tcp(db:3306)/anon_solicitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrate(db)

	db.LogMode(true)
	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&Feedback{})
	db.AutoMigrate(&Questions{})
}
