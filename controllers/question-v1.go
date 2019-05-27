package controllers

import (
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"
	domain "github.com/mike-webster/anon-solicitor/app"
)

func postQuestionV1(c *gin.Context) {
	es, _, _, err := getDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, ErrRetrievingDependencies)

		return
	}

	tok, err := domain.String(c, "tok")
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusUnauthorized)
		setError(c, err, ErrBadToken)

		return
	}

	eventID, _ := strconv.Atoi(c.Param("eventid"))
	event := es.GetEvent(int64(eventID))
	if event == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, err, ErrRetrievingDomainObject)

		return
	}

	// TODO: I need to add the questions to the feedback model
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}
