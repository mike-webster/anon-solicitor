package tokens

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestTokens(t *testing.T) {
	secret := "testsecret"
	jwt := GetJWT(secret, 1, false)
	assert.Equal(t, true, len(jwt) > 0, jwt)
	assert.Equal(t, nil, CheckToken(jwt, secret))
}
