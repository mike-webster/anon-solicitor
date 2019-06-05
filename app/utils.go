package app

import (
	"context"
	"errors"
	"fmt"

	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/anon-solicitor/env"
)

// ContextKey is meant to be used with contexts
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

// DB retrieves the sqlx DB from the gin context
func DB(ctx context.Context) (*sqlx.DB, error) {
	cfg := env.Config()

	db, err := sqlx.Open("mysql", cfg.ConnectionString)
	if err != nil {
		panic(err)
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

func MapStringInterface(ctx *gin.Context, key interface{}) (map[string]interface{}, error) {
	var ret map[string]interface{}
	if ctx == nil {
		return ret, errors.New("provide a gin context in order to retrieve a value")
	}

	ret, ok := ctx.Value(key).(map[string]interface{})
	if !ok {
		return ret, errors.New(fmt.Sprint("couldnt parse map[string]interface from context - ", ctx.Value(key)))
	}

	return ret, nil
}

// String retrieves the expected string value with the given key from the gin context
func String(ctx *gin.Context, key interface{}) (string, error) {
	if ctx == nil {
		return "", errors.New("provide a gin context in order to retrieve a value")
	}

	s, ok := ctx.Value(key).(string)
	if !ok {
		return "", fmt.Errorf("couldnt parse string from context - key: [%v] : value [%v]", key, s)
	}

	return s, nil
}

type DeliveryService interface {
	SendFeedbackEmail(string, string) error
}

var EmailServiceKey ContextKey = "EmailService"
