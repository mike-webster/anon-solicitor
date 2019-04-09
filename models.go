package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

// DBWrapper is a wrapper around the gorm database
type DBWrapper struct {
	db *sqlx.DB
}

// Get will return the open connection if there is one and it
// will try to restablish it if the connection has been closed.
func (w *DBWrapper) Get() *sqlx.DB {
	if w.db != nil {
		fmt.Println("DBWrapper#GET - db found")
		err := w.db.Ping()

		if err != nil {
			fmt.Println("DBWrapper#GET - db error 1: ", err)
			return w.db
		}
	}

	// exactly the same as the built-in
	db, err := sqlx.Open("mysql", "root@tcp(db:3306)/anon_solicitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("DBWrapper#GET - db error 2: ", err)
	}

	// force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		fmt.Println("DBWrapper#GET - db ping: ", err)
		return nil
	}

	w.db = db

	return db
}

// EventPutParams represents the information about an Event that a user can update.
type EventPutParams struct {
	ID          int64     `json:"id" binding:"required"`
	Title       string    `json:"title"`
	Description string    `json:"description" binding:"max=1000"`
	Time        time.Time `json:"scheduled_time"`
	OrganizerID int64
}
