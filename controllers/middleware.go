package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mike-webster/anon-solicitor/email"

	gin "github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/anon-solicitor/data"
	"github.com/mike-webster/anon-solicitor/env"
	"github.com/mike-webster/anon-solicitor/tokens"
	"gopkg.in/go-playground/validator.v8"
)

func setDependencies(ctx context.Context, db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := env.Config()
		es := data.EventService{DB: db}
		fs := data.FeedbackService{DB: db}
		em := email.DeliveryService{
			Host: cfg.SMTPHost,
			Port: cfg.SMTPPort,
			User: cfg.SMTPUser,
			Pass: cfg.SMTPPass,
		}
		// TODO: Fix these context keys. It looks like I'm recreating these amongst pacakges...
		//       Try to just move them into the ENV package as exposed constants?
		c.Set(eventServiceKey.String(), es)
		c.Set(feedbackServiceKey.String(), fs)
		c.Set("EmailService", em)
		c.Next()
	}
}

func getToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Note: this is probably unncessary if the token is going to be a url param...
		//       I just wanted to do it. :)
		// TODO: test this
		cfg := env.Config()
		token := c.Request.Header.Get("token")
		if len(token) < 1 {
			log.Println("token not found in header, checking cookies")
			cookieToken, _ := c.Cookie("anonauth")
			if len(cookieToken) > 1 {
				log.Println("token found in cookie")

				tok, err := tokens.CheckToken(cookieToken, env.Config().Secret)
				if err != nil {
					log.Println("token invalid - 401 - ", tok, " - ", err)
					c.AbortWithError(http.StatusUnauthorized, err)

					return
				}

				c.Set("tok", tok)
				c.Next()

				return
			}

			log.Println("token not found in header, checking query string")
			token = c.Param("token")
			if len(token) < 1 {
				log.Println("token not found - 401")
				c.AbortWithStatus(http.StatusUnauthorized)

				return
			}

			if len(token) < 1 {
				log.Println("no tok in jwt - 401")
				c.AbortWithError(http.StatusUnauthorized, errors.New("couldn't find token"))

				return
			}

			log.Println(fmt.Sprint("tok: ", token))
			c.Set("tok", token)
			c.Next()

			return
		}

		tok, err := tokens.CheckToken(token, cfg.Secret)
		if err != nil {
			log.Println("token invalid - 401 - ", err)
			c.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		if len(tok) < 1 {
			log.Println("no tok in jwt - 401")
			c.AbortWithError(http.StatusUnauthorized, errors.New("couldn't find token"))

			return
		}

		log.Println(fmt.Sprint("tok: ", tok))
		c.Set("tok", tok)
		c.Next()
	}
}

func setStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		_, controllerErr := c.Get("controllerError")
		if controllerErr {
			status, ok := c.Get(controllerRespStatusKey)
			respStatus := http.StatusInternalServerError
			if ok {
				log.Println(fmt.Sprint("error status found: ", status))
				s, ok := status.(int)
				if !ok {
					c.Error(gin.Error{
						Err:  fmt.Errorf("Error processing resp status as int: %v", respStatus),
						Meta: "middleware.setStatus",
					})
				} else {
					respStatus = s
					fmt.Println(fmt.Sprint("updated status: ", respStatus))
				}
			} else {
				log.Println("responseStatus not found - defaulting to 500")
			}

			ret := map[string]string{}

			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypeBind:
					helpful := e.Err.(validator.ValidationErrors)
					for _, err := range helpful {
						ret[err.Field] = ValidationErrorToText(err)
					}
				default:
				}
			}

			c.AbortWithStatusJSON(respStatus, ret)

			return
		}

		// TODO: Should we do something to check for the errors
		// in C.Errors? At this point, you need to remember to
		// do `c.Set("controllerError", true)` in order to
		// have the status set as 500, so if I were to add an
		// error but not set the context value it would render
		// as a 200... but I think errors would get logged?
	}
}

// Below is borrowed from some very kind stranger.
// sauce: https://github.com/gin-gonic/gin/issues/430#issuecomment-446113460
var (
	ErrorInternalError = errors.New("whoops something went wrong")
)

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	return strings.ToLower(str)
}

func Split(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	for index, word := range entries {
		if index == 0 {
			entries[index] = UcFirst(word)
		} else {
			entries[index] = LcFirst(word)
		}
	}
	justString := strings.Join(entries, " ")
	return justString
}

func ValidationErrorToText(e *validator.FieldError) string {
	word := Split(e.Field)

	switch e.Tag {
	case "required":
		return fmt.Sprintf("%s is required", word)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", word, e.Param)
	case "min":
		return fmt.Sprintf("%s must be longer than %s", word, e.Param)
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", word, e.Param)
	}
	return fmt.Sprintf("%s is not valid", word)
}

// This method collects all errors and submits them to Rollbar
// func Errors() gin.HandlerFunc {

// 	return func(c *gin.Context) {
// 		c.Next()
// 		// Only run if there are some errors to handle
// 		if len(c.Errors) > 0 {
// 			for _, e := range c.Errors {
// 				// Find out what type of error it is
// 				switch e.Type {
// 				case gin.ErrorTypePublic:
// 					// Only output public errors if nothing has been written yet
// 					if !c.Writer.Written() {
// 						c.JSON(c.Writer.Status(), gin.H{"Error": e.Error()})
// 					}
// 				case gin.ErrorTypeBind:
// 					errs := e.Err.(validator.ValidationErrors)
// 					list := make(map[string]string)
// 					for _,err := range errs {
// 						list[err.Field] = ValidationErrorToText(err)
// 					}

// 					// Make sure we maintain the preset response status
// 					status := http.StatusBadRequest
// 					if c.Writer.Status() != http.StatusOK {
// 						status = c.Writer.Status()
// 					}
// 					c.JSON(status, gin.H{"Errors": list})

// 				default:
// 					// Log all other errors
// 					rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
// 				}

// 			}
// 			// If there was no public or bind error, display default 500 message
// 			if !c.Writer.Written() {
// 				c.JSON(http.StatusInternalServerError, gin.H{"Error": ErrorInternalError.Error()})
// 			}
// 		}
// 	}
// }
