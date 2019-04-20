package main

import (
	"context"
	"flag"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mike-webster/anon-solicitor/controllers"
)

// TODO: better logging?

func main() {
	log.Print("Sleeping to allow db setup...")
	time.Sleep(3 * time.Second)

	ctx := moveFlagsToContext()

	controllers.StartServer(ctx)
}

func moveFlagsToContext() context.Context {
	ctx := context.Background()

	dropTables := flag.Bool("drop", false, "should we drop the existing tables")
	flag.Parse()
	log.Printf("-- DropTables: %v", *dropTables)
	ctx = context.WithValue(ctx, "DropTables", *dropTables)

	return ctx
}
