package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
	domain "github.com/mike-webster/anon-solicitor/app"
)

func postQuestionV1(c *gin.Context) {
	es, _, _, err := GetDependencies(c)
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
		Content:   postQuestion.Content,
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

func getQuestionV1(c *gin.Context) {
	es, fs, _, err := GetDependencies(c)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, ErrRetrievingDependencies)

		return
	}

	tokPayload, err := domain.MapStringInterface(c, "tok")
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusUnauthorized)
		setError(c, err, ErrBadToken)

		return
	}

	eid, _ := tokPayload["eid"].(float64)
	event := es.GetEvent(int64(eid))
	if event == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNotFound)
		setError(c, err, "couldnt find event")

		return
	}

	tok, _ := tokPayload["tok"].(string)
	questions := fs.GetQuestionsForTok(tok)
	if questions == nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusNoContent)
		setError(c, err, "no_questions_for_tok")

		return
	}

	display := []domain.QuestionDisplay{}
	for _, q := range *questions {
		display = append(display,
			domain.QuestionDisplay{
				QuestionID: q.ID,
				EventID:    q.EventID,
				Content:    q.Content,
				Answers:    q.Answers(),
			})
	}

	c.HTML(http.StatusOK, "question.html", gin.H{
		"title":       fmt.Sprint("AnonSolicitor - ", event.Title),
		"headline":    event.Title,
		"description": event.Description,
		"questions":   display,
	})
}

func postQuestionAnswerV1(c *gin.Context) {
	es, fs, _, err := GetDependencies(c)
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
		setError(c, errors.New(fmt.Sprint("couldnt find question: ", questionID)), ErrRetrievingDomainObject)

		return
	}

	token, _ := payload["tok"].(string)
	if !es.CanUserAnswerQuestion(question.ID, token) {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		setError(c, errors.New("user cannot answer question"), "err_forbidden")

		return
	}

	feedback, err := fs.GetFeedbackByTok(token)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusForbidden)
		setError(c, fmt.Errorf("Feedback error: %v", err), "err_forbidden")

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
		Content:    postAnswer.Content,
		QuestionID: question.ID,
	}

	err = es.AddAnswer(&newAnswer, feedback.ID)
	if err != nil {
		c.Set(controllerErrorKey, true)
		c.Set(controllerRespStatusKey, http.StatusInternalServerError)
		setError(c, err, "err_persisting_answer")

		return
	}

	// TODO: I need to add the questions to the feedback model
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}
