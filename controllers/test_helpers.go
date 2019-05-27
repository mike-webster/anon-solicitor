package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func performRequest(r http.Handler, method string, path string, body *[]byte, headers map[string]string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		reader := strings.NewReader(string(*body))
		log.Println("Sending test body")
		req, _ = http.NewRequest(method, path, reader)
	} else {
		log.Println("Sending empty test body")
		req, _ = http.NewRequest(method, path, nil)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
