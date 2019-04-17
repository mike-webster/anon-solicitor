package main

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/anon-solicitor/data"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mike-webster/anon-solicitor/controllers"
)

func main() {
	ctx := context.Background()
	log.Print("Sleeping to allow db setup...")
	time.Sleep(3 * time.Second)

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

	controllers.StartServer(ctx, es)
}
