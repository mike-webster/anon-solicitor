package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/anon-solicitor/app"
)

func TestTestSetup(t *testing.T) {
	t.Run("TestGetTestSetup", func(t *testing.T) {
		headers := getTestHeaders()
		opts := app.TestServiceOptions{}
		deps := app.MockSearchDependencies(opts)
		t.Run("WithCookie", func(t *testing.T) {
			headers["Cookie"] = fmt.Sprint("anonauth=", "notarealtoken")
			t.Run("WhenErrorOccursCheckingToken", func(t *testing.T) {
				t.Run("ItReturnsDefaultData", func(t *testing.T) {
					r := setupTestRouter(deps, false)
					req := performRequest(r, "GET", "/testsetup", nil, headers)
					assert.Equal(t, 200, req.Code, req.Body.String())
					fragileTestPattern := `<input type="text" name="token" value=".*" style="width:50%;min-width:50px;"></p>`
					matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
					assert.Equal(t, nil, err)
					assert.Equal(t, true, matched)
				})
			})
			t.Run("WhenTheRoleIsNotAString", func(t *testing.T) {
				t.Run("ItReturnsDefaultData", func(t *testing.T) {
					r := setupTestRouter(deps, false)
					tokv, _ := uuid.NewV4()
					payload := map[string]interface{}{
						"eid":  30,
						"role": 1234.432,
						"tok":  tokv.String(),
					}
					tok := getTestTok(&payload)
					headers["Cookie"] = fmt.Sprintf("anonauth=%v", tok)
					req := performRequest(r, "GET", "/testsetup", nil, headers)
					assert.Equal(t, 200, req.Code, req.Body.String())
					fragileTestPattern := `<p>Role:&nbsp;<input type="text" name="role" value="audience" style="min-width:50px;"></p>`
					matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
					assert.Equal(t, nil, err)
					assert.Equal(t, true, matched)
				})
			})
			t.Run("WhenTheEIDIsNotAnInt64", func(t *testing.T) {
				t.Run("ItReturnsDefaultData", func(t *testing.T) {
					r := setupTestRouter(deps, false)
					tokv, _ := uuid.NewV4()
					payload := map[string]interface{}{
						"eid":  tokv.String(),
						"role": RoleAudience,
						"tok":  tokv.String(),
					}
					tok := getTestTok(&payload)
					headers["Cookie"] = fmt.Sprintf("anonauth=%v", tok)
					req := performRequest(r, "GET", "/testsetup", nil, headers)
					assert.Equal(t, 200, req.Code, req.Body.String())
					fragileTestPattern := `<p>Event ID:&nbsp;<input type="text" name="eventid" value="1"></p>`
					matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
					assert.Equal(t, nil, err)
					assert.Equal(t, true, matched)
				})
			})
			t.Run("WhenProvidedAValidToken", func(t *testing.T) {
				r := setupTestRouter(deps, false)
				tokv, _ := uuid.NewV4()
				expectedEID := 501
				expectedAudience := RoleOwner
				expectedTok := tokv.String()
				pl := map[string]interface{}{
					"eid":  expectedEID,
					"role": expectedAudience,
					"tok":  expectedTok,
				}
				tok := getTestTok(&pl)
				headerss := getTestHeaders()
				headerss["Cookie"] = fmt.Sprintf("anonauth=%v", tok)
				req := performRequest(r, "GET", "/testsetup", nil, headerss)
				t.Run("ItReturns200", func(t *testing.T) {
					assert.Equal(t, http.StatusOK, req.Code, req.Body.String())
				})
				t.Run("ItPopulatesTheEventID", func(t *testing.T) {
					fragileTestPattern := fmt.Sprintf(`<p>Event ID:&nbsp;<input type="text" name="eventid" value="%v"></p>`, expectedEID)
					matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
					assert.Equal(t, nil, err)
					assert.Equal(t, true, matched, expectedEID)
				})
				t.Run("ItPopulatesTheRole", func(t *testing.T) {
					fragileTestPattern := fmt.Sprintf(`<p>Role:&nbsp;<input type="text" name="role" value="%v" style="min-width:50px;"></p>`, expectedAudience)
					matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
					assert.Equal(t, nil, err)
					assert.Equal(t, true, matched)
				})
				t.Run("ItPopulatesTheToken", func(t *testing.T) {
					fragileTestPattern := fmt.Sprintf(`<p>Token:&nbsp;<input type="text" name="token" value="%v"`, expectedTok)
					matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
					assert.Equal(t, nil, err)
					assert.Equal(t, expectedTok, tokv.String(), fmt.Sprintf("%v != %v", expectedTok, tokv.String()))
					assert.Equal(t, true, matched, expectedTok, "\n", req.Body.String())
				})
			})
		})
	})

	t.Run("TestPostTestSetup", func(t *testing.T) {
		t.Run("WhenBindingErrorOccurrs", func(t *testing.T) {
			t.Run("ItReturns500", func(t *testing.T) {

			})
		})
		t.Run("WhenTokenIsntProvided", func(t *testing.T) {
			t.Run("ItGeneratesANewOne", func(t *testing.T) {

			})
		})
		t.Run("WhenEventIdIsntProvided", func(t *testing.T) {
			t.Run("ItIsSetTo1", func(t *testing.T) {

			})
		})
		t.Run("WhenTheFormIsValid", func(t *testing.T) {
			headers := getTestHeaders()
			opts := app.TestServiceOptions{}
			deps := app.MockSearchDependencies(opts)
			r := setupTestRouter(deps, false)
			tokv, _ := uuid.NewV4()
			expectedEID := int64(501)
			expectedAudience := RoleOwner
			expectedTok := tokv.String()
			pl := map[string]interface{}{
				"eid":  expectedEID,
				"role": expectedAudience,
				"tok":  expectedTok,
			}
			form := testSetup{
				Token:   expectedTok,
				EventID: expectedEID,
				Role:    expectedAudience,
			}
			formBytes, err := json.Marshal(&form)
			assert.Equal(t, nil, err)

			tok := getTestTok(&pl)
			req := performRequest(r, "POST", "/testsetup", &formBytes, headers)
			t.Run("ItReturns200", func(t *testing.T) {
				assert.Equal(t, http.StatusOK, req.Code, req.Body.String())
			})
			t.Run("ItSetsTheCookie", func(t *testing.T) {
				cookie := req.HeaderMap["Set-Cookie"][0]

				fragileTestPattern := fmt.Sprintf(`anonauth=%v; Path=/; Domain=localhost; Max-Age=3600`, url.QueryEscape(tok))
				matched, err := regexp.MatchString(fragileTestPattern, cookie)
				assert.Equal(t, nil, err)
				assert.Equal(t, true, matched, cookie, "\n", url.QueryEscape(tok))
			})
			t.Run("ItPopulatesTheEventID", func(t *testing.T) {
				fragileTestPattern := fmt.Sprintf(`<p>Event ID:&nbsp;<input type="text" name="eventid" value="%v"></p>`, expectedEID)
				matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
				assert.Equal(t, nil, err)
				assert.Equal(t, true, matched, expectedEID)
			})
			t.Run("ItPopulatesTheRole", func(t *testing.T) {
				fragileTestPattern := fmt.Sprintf(`<p>Role:&nbsp;<input type="text" name="role" value="%v" style="min-width:50px;"></p>`, expectedAudience)
				matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
				assert.Equal(t, nil, err)
				assert.Equal(t, true, matched)
			})
			t.Run("ItPopulatesTheToken", func(t *testing.T) {
				fragileTestPattern := fmt.Sprintf(`<p>Token:&nbsp;<input type="text" name="token" value="%v"`, expectedTok)
				matched, err := regexp.MatchString(fragileTestPattern, req.Body.String())
				assert.Equal(t, nil, err)
				assert.Equal(t, expectedTok, tokv.String(), fmt.Sprintf("%v != %v", expectedTok, tokv.String()))
				assert.Equal(t, true, matched, expectedTok, "\n", req.Body.String())
			})
		})
	})
}
