package controllers

import (
	"testing"

	"github.com/bmizerany/assert"
	anon "github.com/mike-webster/anon-solicitor"
)

func TestContextKeys(t *testing.T) {
	val := "test"
	key := anon.ContextKey(val)
	assert.Equal(t, val, key.String())
}
