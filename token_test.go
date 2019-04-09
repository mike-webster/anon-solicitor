package main

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestTokens(t *testing.T) {
	jwt := GetJWT()
	assert.Equal(t, true, len(jwt) > 0, jwt)
	assert.Equal(t, nil, CheckToken(jwt))
}
