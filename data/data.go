package data

import (
	"context"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// CreateTables will create the necessary tables for the application
// if they do not exist.
// Note: if the context has a true value for "DropValues", this flow
// will drop the existing tables if they exist in order to start with
// a fresh databse.
func CreateTables(ctx context.Context, db *sqlx.DB) error {
	dt, ok := ctx.Value("DropTables").(bool)
	if !ok {
		return errors.New("couldnt parse DropTables from context")
	}
	if dt {
		dropTables(ctx)
	}
	fmt.Println("-- Creating tables")

	eventSchema := `CREATE TABLE IF NOT EXISTS events (
		id INT AUTO_INCREMENT,
		title NVARCHAR(200) NOT NULL,
		description NVARCHAR(5000),
		time DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME,
		deleted_at DATETIME,
		PRIMARY KEY (id)
	);`
	err := createTable(ctx, db, "events", eventSchema)
	if err != nil {
		return err
	}

	// feedback
	feedbackSchema := `CREATE TABLE IF NOT EXISTS feedback (
		id INT AUTO_INCREMENT,
		content NVARCHAR(5000) NOT NULL,
		tok NVARCHAR(5000) NOT NULL,
		event_id INT NOT NULL,
		absent BOOLEAN NOT NULL DEFAULT FALSE,
		PRIMARY KEY(id)
	);`
	err = createTable(ctx, db, "feedback", feedbackSchema)
	if err != nil {
		return err
	}

	// questions
	questionSchema := `CREATE TABLE IF NOT EXISTS questions (
		id INT AUTO_INCREMENT,
		event_id INT NOT NULL,
		content NVARCHAR(5000) NOT NULL,
		answers NVARCHAR(5000) NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME,
		deleted_at DATETIME,
		PRIMARY KEY (id)
	);`
	err = createTable(ctx, db, "questions", questionSchema)
	if err != nil {
		return err
	}

	return nil
}

func createTable(ctx context.Context, db *sqlx.DB, tableName string, tableSchema string) error {
	_, err := db.Exec(tableSchema)
	if err != nil {
		return err
	}

	return nil
}

func dropTables(ctx context.Context) error {
	log.Println("-- Deleting tables")
	queries := []string{
		"DROP TABLE IF EXISTS feedback;",
		"DROP TABLE IF EXISTS events;",
		"DROP TABLE IF EXISTS questions;",
	}

	db, err := sqlx.Open("mysql", "root@tcp(db:3306)/anon_solicitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("DB - db error 2: ", err)

		return nil
	}
	defer db.Close()

	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
