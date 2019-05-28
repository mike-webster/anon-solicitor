package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mike-webster/anon-solicitor/email"

	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/anon-solicitor/data"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
)

func setDependencies(ctx context.Context, db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := env.Config()
		es := data.EventService{DB: db}
		fs := data.FeedbackService{DB: db}
		em := email.DeliveryService{
			Host: cfg.SMTPHost,
			Port: cfg.SMTPPort,
			User: cfg.SMTPUser,
			Pass: cfg.SMTPPass,
		}
		// TODO: Fix these context keys. It looks like I'm recreating these amongst pacakges...
		//       Try to just move them into the ENV package as exposed constants?
		c.Set(eventServiceKey.String(), es)
		c.Set(feedbackServiceKey.String(), fs)
		c.Set("EmailService", em)
		c.Next()
	}
}

func getToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Note: this is probably unncessary if the token is going to be a url param...
		//       I just wanted to do it. :)
		// TODO: test this
		cfg := env.Config()
		token := c.Request.Header.Get("token")
		if len(token) < 1 {
			log.Println("token not found in header, checking query string")

			token := c.Param("token")
			if len(token) < 1 {
				log.Println("token not found - 401")
				c.AbortWithStatus(http.StatusUnauthorized)

				return
			}
		}

		tok, err := tokens.CheckToken(token, cfg.Secret)
		if err != nil {
			log.Println("token invalid - 401 - ", err)
			c.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		if len(tok) < 1 {
			log.Println("no tok in jwt - 401")
			c.AbortWithError(http.StatusUnauthorized, errors.New("couldn't find token"))

			return
		}

		log.Println(fmt.Sprint("tok: ", tok))
		c.Set("tok", tok)
		c.Next()
	}
}

func setStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		_, controllerErr := c.Get("controllerError")
		if controllerErr {
			status, ok := c.Get(controllerRespStatusKey)
			respStatus := http.StatusInternalServerError
			if ok {
				respStatus, ok = status.(int)
				if !ok {
					c.Error(gin.Error{
						Err:  fmt.Errorf("Error processing resp status as int: %v", respStatus),
						Meta: "middleware.setStatus",
					})
				}
			} else {
				log.Print("responseStatus not found - defaulting to 500")
			}
			c.AbortWithStatusJSON(respStatus, gin.H{"msg": "sorry - we encountered an error and we're working on it!", "errors": c.Errors})

			return
		}

		// TODO: Should we do something to check for the errors
		// in C.Errors? At this point, you need to remember to
		// do `c.Set("controllerError", true)` in order to
		// have the status set as 500, so if I were to add an
		// error but not set the context value it would render
		// as a 200... but I think errors would get logged?
	}
}
