package app

import (
	"errors"
	"fmt"

	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ContextKey is meant to be used with contexts
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

// DB retrieves the sqlx DB from the gin context
func DB(ctx *gin.Context) (*sqlx.DB, error) {
	if ctx == nil {
		return nil, errors.New("provide a gin context in order to retrieve the database")
	}

	db, ok := ctx.Value("DB").(*sqlx.DB)
	if !ok {
		return nil, errors.New("couldnt parse db from context")
	}

	return db, nil
}

// Bool retrieves the expected bool value with the given key from the gin context
func Bool(ctx *gin.Context, key interface{}) (*bool, error) {
	if ctx == nil {
		return nil, errors.New("provide a gin context in order to retrieve a value")
	}

	b, ok := ctx.Value(key).(bool)
	if !ok {
		return nil, fmt.Errorf("couldnt parse bool from context for key: %v", key)
	}

	return &b, nil
}

// String retrieves the expected string value with the given key from the gin context
func String(ctx *gin.Context, key interface{}) (string, error) {
	if ctx == nil {
		return "", errors.New("provide a gin context in order to retrieve a value")
	}

	s, ok := ctx.Value(key).(string)
	if !ok {
		return "", errors.New("couldnt parse string from context")
	}

	return s, nil
}

type DeliveryService interface {
	SendFeedbackEmail(string, string) error
}

var EmailServiceKey ContextKey = "EmailService"
