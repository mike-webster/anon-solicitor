package tokens

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/mike-webster/anon-solicitor/env"
)

var cfg = env.Config()

func TestTokens(t *testing.T) {
	tok := "thisisatestoneusetoken"
	payload := map[string]interface{}{
		"tok": tok,
		"exp": time.Now().UTC().Add(30 * time.Minute).Unix(),
		"iss": "anon-test",
	}
	jwt := GetJWT(cfg.Secret, payload)
	assert.Equal(t, true, len(jwt) > 0, jwt)
	ret, err := CheckToken(jwt, cfg.Secret)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, tok, ret)
}
