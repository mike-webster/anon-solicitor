package tokens

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/mike-webster/anon-solicitor/env"
)

var cfg = env.Config()

func TestTokens(t *testing.T) {
	tok := "thisisatestoneusetoken"
	jwt := GetJWT(cfg.Secret, tok)
	assert.Equal(t, true, len(jwt) > 0, jwt)
	ret, err := CheckToken(jwt, cfg.Secret)
	assert.Equal(t, nil, err)
	assert.Equal(t, tok, ret)
}
