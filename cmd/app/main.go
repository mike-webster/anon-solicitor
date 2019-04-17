package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/anon-solicitor/data"
	"github.com/mike-webster/anon-solicitor/env"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mike-webster/anon-solicitor/controllers"
)

func main() {
	log.Print("Sleeping to allow db setup...")
	time.Sleep(3 * time.Second)

	var _ = env.Config()

	ctx := moveFlagsToContext()

	err := data.CreateTables(ctx)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Open("mysql", "root@tcp(db:3306)/anon_solicitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	es := data.EventService{DB: db}
	fs := data.FeedbackService{DB: db}

	controllers.StartServer(ctx, es, fs)
}

func moveFlagsToContext() context.Context {
	ctx := context.Background()

	dropTables := flag.Bool("drop", false, "should we drop the existing tables")
	flag.Parse()
	log.Printf("-- DropTables: %v", *dropTables)
	ctx = context.WithValue(ctx, "DropTables", *dropTables)

	return ctx
}
