package app

import (
	"testing"

	"github.com/bmizerany/assert"
	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func TestDB(t *testing.T) {
	t.Run("NilContextReturnsError", func(t *testing.T) {
		_, err := DB(nil)
		assert.NotEqual(t, nil, err)
	})

	t.Run("InvalidValueForKey", func(t *testing.T) {
		ctx := gin.Context{}
		ctx.Set("DB", "invalid")
		_, err := DB(&ctx)
		assert.NotEqual(t, nil, err)
	})

	t.Run("Valid", func(t *testing.T) {
		ctx := gin.Context{}
		db := sqlx.DB{}
		ctx.Set("DB", &db)
		_, err := DB(&ctx)
		assert.Equal(t, nil, err)
	})
}

func TestBool(t *testing.T) {
	t.Run("NilContextReturnsError", func(t *testing.T) {
		_, err := Bool(nil, "")
		assert.NotEqual(t, nil, err)
	})

	t.Run("InvalidValueForKey", func(t *testing.T) {
		ctx := gin.Context{}
		key := "test"
		ctx.Set(key, "testing")
		_, err := Bool(&ctx, key)
		assert.NotEqual(t, nil, err)
	})

	t.Run("Valid", func(t *testing.T) {
		ctx := gin.Context{}
		key := "test"
		ctx.Set(key, true)
		v, err := Bool(&ctx, key)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, *v)
	})
}

func TestString(t *testing.T) {
	t.Run("NilContextReturnsError", func(t *testing.T) {
		_, err := String(nil, "")
		assert.NotEqual(t, nil, err)
	})

	t.Run("InvalidValueForKey", func(t *testing.T) {
		ctx := gin.Context{}
		key := "test"
		ctx.Set(key, 23456)
		_, err := String(&ctx, key)
		assert.NotEqual(t, nil, err)
	})

	t.Run("Valid", func(t *testing.T) {
		ctx := gin.Context{}
		key := "test"
		val := "testvalue"
		ctx.Set(key, val)
		v, err := String(&ctx, key)
		assert.Equal(t, nil, err)
		assert.Equal(t, val, v)
	})
}
