package main

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

// DBWrapper is a wrapper around the gorm database
type DBWrapper struct {
	db *gorm.DB
}

// Get will return the open connection if there is one and it
// will try to restablish it if the connection has been closed.
func (w *DBWrapper) Get() *gorm.DB {
	connection := os.Getenv("DB_USER") +
		":" +
		os.Getenv("DB_PASS") +
		"@tcp(" +
		os.Getenv("DB_HOST") +
		":" +
		os.Getenv("DB_PORT") +
		")/anon_solicitor?charset=utf8&parseTime=True&loc=Local"
	if w.db == nil {
		db, err := gorm.Open("mysql", connection)
		if err != nil {
			log.Printf("db error: %v", err)
			return nil
		}

		w.db = db
	}

	err := w.db.DB().Ping()
	if err != nil {
		log.Printf("DBWrapper.Get Error: %v", err)

		db, err := gorm.Open("mysql", connection)
		if err != nil {
			log.Printf("db last ditch error: %v", err)
			return nil
		}

		w.db = db
	}

	return w.db
}

// User represents an application user.
type User struct {
	gorm.Model
	Name  string
	Email string
}

// Feedback represents a user's opinion on an event.
type Feedback struct {
	gorm.Model
	Content string
}

// Event represents a situation about which a user would like anonymous feedback.
type Event struct {
	gorm.Model
	Name               string      `json:"name" binding:"required,max=100"`
	Description        string      `json:"description" binding:"required,max=1000"`
	Time               time.Time   `json:"scheduled_time" binding:"required"`
	OrganizingUser     User        `binding:"-"`
	OrganizerQuestions []Questions `binding:"-"`
	Feedback           []Feedback  `binding:"-"`
}

type Questions struct {
	gorm.Model
	Content string
	Answers string
}
