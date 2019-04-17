package data

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var _db *sqlx.DB

func DB() *sqlx.DB {
	if _db == nil {
		db, err := sqlx.Open("mysql", "root@tcp(db:3306)/anon_solicitor?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Println("DB - db error 2: ", err)

			return nil
		}

		_db = db
	}

	err := _db.Ping()
	if err != nil {
		fmt.Println("DB - ping err: ", err)

		return nil
	}

	return _db
}

func DropTables(ctx context.Context) error {
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

func CreateTables(ctx context.Context) error {
	dropTables, _ := ctx.Value("DropTables").(bool)
	if dropTables {
		DropTables(ctx)
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
	err := createTable(ctx, "events", eventSchema)
	if err != nil {
		return err
	}

	// feedback
	feedbackSchema := `CREATE TABLE IF NOT EXISTS feedback (
		id INT AUTO_INCREMENT,
		content NVARCHAR(5000) NOT NULL,
		tok NVARCHAR(5000) NOT NULL,
		event_id INT NOT NULL,
		PRIMARY KEY(id)
	);`
	err = createTable(ctx, "feedback", feedbackSchema)
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
	err = createTable(ctx, "questions", questionSchema)
	if err != nil {
		return err
	}

	return nil
}

func createTable(ctx context.Context, tableName string, tableSchema string) error {
	db, err := sqlx.Open("mysql", "root@tcp(db:3306)/anon_solicitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("DB - db error 2: ", err)

		return nil
	}
	defer db.Close()

	_, err = db.Exec(tableSchema)
	if err != nil {
		return err
	}

	return nil
}
