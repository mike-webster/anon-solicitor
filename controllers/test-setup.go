package controllers

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func getTestSetup(c *gin.Context) {
	tok, _ := uuid.NewV4()
	c.HTML(http.StatusOK, "testsetup.html", gin.H{"tok": tok})
}

func postTestSetup(c *gin.Context) {

}
