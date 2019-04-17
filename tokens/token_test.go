package tokens

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestTokens(t *testing.T) {
	secret := "testsecret"
	tok := "thisisatestoneusetoken"
	jwt := GetJWT(secret, tok)
	assert.Equal(t, true, len(jwt) > 0, jwt)
	ret, err := CheckToken(jwt, secret)
	assert.Equal(t, nil, err)
	assert.Equal(t, tok, ret)
}
