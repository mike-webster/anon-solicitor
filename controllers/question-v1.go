package controllers

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

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

	payload, err := domain.MapStringInterface(c, "tok")
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusUnauthorized)
		setError(c, err, ErrBadToken)

		return
	}

	if payload["role"] != RoleOwner {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		log.Println("~Forbidden #1")
		setError(c, err, ErrNotAllowed)

		return
	}

	eventID, _ := strconv.Atoi(c.Param("eventid"))
	log.Println("eid: ", payload["eid"], " - ", reflect.TypeOf(payload["eid"]))
	eid, _ := payload["eid"].(float64)
	if int64(eid) != int64(eventID) {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		log.Println("~Forbidden #2 : ", eid, " != ", eventID)
		setError(c, errors.New("user trying to manipulate feedback for unintended event"), "error_mismatched_ids")

		return
	}

	event := es.GetEvent(int64(eventID))
	if event == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, errors.New("couldnt find event"), ErrRetrievingDomainObject)

		return
	}
	postQuestion := domain.QuestionPostParams{}
	err = c.Bind(&postQuestion)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusBadRequest)
		setError(c, err, ErrValidation)
		log.Printf("Error binding object: %v", err)

		return
	}

	log.Printf("posted question: %v", postQuestion)

	newQuestion := domain.Question{
		EventID:   event.ID,
		Title:     postQuestion.Title,
		Content:   postQuestion.Content,
		Answers:   postQuestion.Answers,
		CreatedAt: time.Now().UTC(),
	}

	err = es.AddQuestion(&newQuestion)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, "err_persisting_question")

		return
	}

	// TODO: I need to add the questions to the feedback model
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func postQuestionAnswerV1(c *gin.Context) {
	es, _, _, err := getDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, ErrRetrievingDependencies)

		return
	}

	payload, err := domain.MapStringInterface(c, "tok")
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusUnauthorized)
		setError(c, err, ErrBadToken)

		return
	}

	if payload["role"] != RoleOwner && payload["role"] != RoleAudience {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		log.Println("~Forbidden #1")
		setError(c, err, ErrNotAllowed)

		return
	}

	eventID, _ := strconv.Atoi(c.Param("eventid"))
	eid, _ := payload["eid"].(float64)
	if int64(eid) != int64(eventID) {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		log.Println("~Forbidden #2 : ", eid, " != ", eventID)
		setError(c, errors.New("user trying to manipulate feedback for unintended event"), "error_mismatched_ids")

		return
	}

	event := es.GetEvent(int64(eventID))
	if event == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, errors.New("couldnt find event"), ErrRetrievingDomainObject)

		return
	}

	questionID, _ := strconv.Atoi(c.Param("questionid"))
	qid, _ := payload["eid"].(float64)
	if int64(eid) != int64(eventID) {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		log.Println("~Forbidden #2 : ", qid, " != ", questionID)
		setError(c, errors.New("user trying to manipulate feedback for unintended event"), "error_mismatched_ids")

		return
	}

	question := es.GetQuestion(int64(questionID))
	if question == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, errors.New("couldnt find question"), ErrRetrievingDomainObject)

		return
	}

	token, _ := payload["tok"].(string)
	if !es.CanUserAnswerQuestion(question.ID, token) {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		setError(c, errors.New("user cannot answer question"), "err_forbidden")

		return
	}

	postAnswer := domain.AnswerPostParams{}
	err = c.Bind(&postAnswer)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusBadRequest)
		setError(c, err, ErrValidation)
		log.Printf("Error binding object: %v", err)

		return
	}

	log.Printf("posted answer: %v", postAnswer)

	newAnswer := domain.Answer{
		EventID: event.ID,
		Content: postAnswer.Content,
		Token:   token,
	}

	err = es.AddAnswer(&newAnswer)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, "err_persisting_answer")

		return
	}

	// TODO: I need to add the questions to the feedback model
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}
