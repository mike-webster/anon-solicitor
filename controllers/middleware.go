package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/anon-solicitor/data"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
)

func setDependencies(ctx context.Context, db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		es := data.EventService{DB: db}
		fs := data.FeedbackService{DB: db}
		c.Set(eventServiceKey.String(), es)
		c.Set(feedbackServiceKey.String(), fs)
		c.Next()
	}
}

func getToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Note: this is probably unncessary if the token is going to be a url param...
		//       I just wanted to do it. :)
		// TODO: test this
		cfg := env.Config()
		token := c.Param("token")
		if len(token) < 1 {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		tok, err := tokens.CheckToken(token, cfg.Secret)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		if len(tok) < 1 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("couldn't find token"))

			return
		}

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
