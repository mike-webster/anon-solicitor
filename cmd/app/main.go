package main

import (
	"context"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mike-webster/anon-solicitor/controllers"
)

func main() {
	ctx := context.Background()
	log.Print("Sleeping to allow db setup...")
	time.Sleep(10 * time.Second)

	controllers.StartServer(ctx)

	// wrap := main.DBWrapper{}
	// wrap.Get()

	// ctx = context.WithValue(ctx, "db", wrap)
	// err := data.CreateTables(ctx)
	// if err != nil {
	// 	panic(err)
	// }
}

// func mwAttachDB(c *gin.Context) {
// 	wrap := DBWrapper{}
// 	wrap.Get()
// 	c.Set("db", wrap)

// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("Middleware caught a panic", r)
// 		}
// 	}()

// 	c.Next()
// }
