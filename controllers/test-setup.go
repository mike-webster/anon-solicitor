package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	gin "github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
)

type testSetup struct {
	Token   string `json:"tok" form:"token"`
	EventID int64  `json:"eid" form:"eventid"`
	Role    string `json:"role" form:"role"`
}

func (ts *testSetup) getPayload() map[string]interface{} {
	return map[string]interface{}{
		"tok":  ts.Token,
		"eid":  ts.EventID,
		"role": ts.Role,
	}
}

func getTestSetup(c *gin.Context) {
	tok, _ := uuid.NewV4()
	t := tok.String()
	payload := map[string]interface{}{}
	token, _ := c.Cookie("anonauth")
	if len(token) > 1 {
		pl, err := tokens.CheckToken(token, env.Config().Secret)
		if err != nil {
			log.Println("err checking token: ", err)
		}
		log.Println("found payload: ", pl)
		payload = pl
	} else {
		log.Println("no token found, generating new: ", t)
	}

	form := testSetup{
		Role:    RoleAudience,
		EventID: 1,
	}
	log.Println("payload: ", payload)
	pr, ok := payload["role"].(string)
	if ok {
		log.Println("found role in payload: ", pr)
		form.Role = pr
	} else {
		log.Println("no role found in payload: ", reflect.TypeOf(payload["role"]))
	}

	eid, ok := payload["eid"].(float64)
	if ok {
		log.Println("found eid in payload: ", eid)
		form.EventID = int64(eid)
	} else {
		log.Println("no eid found in payload: ", reflect.TypeOf(payload["eid"]))
	}

	// TODO: populate token, you dumbass
	tok2, ok := payload["tok"].(string)
	if ok {
		log.Println("found tok in payload: ", tok2)
		form.Token = tok2
	} else {
		log.Println("no token found in payload: ", reflect.TypeOf(payload["tok"]))
		form.Token = t
	}

	bytes, _ := json.Marshal(&form)
	log.Println("Form: ", string(bytes))

	c.HTML(http.StatusOK,
		"testsetup.html",
		gin.H{
			"form": form,
		})
}

func postTestSetup(c *gin.Context) {
	var form testSetup
	err := c.Bind(&form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})

		return
	}

	if len(form.Token) < 1 {
		id, _ := uuid.NewV4()
		form.Token = id.String()
	}
	if form.EventID < 1 {
		form.EventID = 1
	}
	if len(form.Role) < 1 {
		form.Role = RoleAudience
	}

	c.SetCookie("anonauth",
		tokens.GetJWT(env.Config().Secret, form.getPayload()),
		3600,
		"/",
		"localhost",
		false,
		false)

	c.HTML(http.StatusOK,
		"testsetup.html",
		gin.H{
			"form": form,
		})
}
