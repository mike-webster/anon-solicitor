package main

import (
	"context"
	"flag"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mike-webster/anon-solicitor/controllers"
	"github.com/mike-webster/anon-solicitor/env"
)

// TODO: better logging?

func main() {
	log.Print("Sleeping to allow db setup...")
	time.Sleep(3 * time.Second)

	ctx := moveFlagsToContext()
	cfg := env.Config()
	r := controllers.GetRouter(ctx)

	r.Run(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port))
}

func moveFlagsToContext() context.Context {
	ctx := context.Background()

	dropTables := flag.Bool("drop", false, "should we drop the existing tables")
	createTables := flag.Bool("create", false, "should we create the existing tables")

	flag.Parse()
	log.Printf("-- DropTables: %v", *dropTables)
	log.Printf("-- CreateTables: %v", *createTables)
	ctx = context.WithValue(ctx, "DropTables", *dropTables)
	ctx = context.WithValue(ctx, "CreateTables", *createTables)

	return ctx
}
