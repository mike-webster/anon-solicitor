package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func performRequest(r http.Handler, method string, path string, body *[]byte) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		reader := strings.NewReader(string(*body))
		log.Println("Sending test body")
		req, _ = http.NewRequest(method, path, reader)
	} else {
		log.Println("Sending empty test body")
		req, _ = http.NewRequest(method, path, nil)
	}
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
