package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {
	Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := NewCorsHandler(Handler)
	response := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/", nil)
	middleware.ServeHTTP(response, req)
	result := response.Result()

	origin := result.Header.Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("Wrong Access Control Allow Origin Handler: expected %s but got %s", origin, "*")
	}

	allowedMethods := "GET, PUT, POST, PATCH, DELETE"
	methods := result.Header.Get("Access-Control-Allow-Methods")
	if methods != allowedMethods {
		t.Errorf("Wrong Access Control Allow Methods Handler: expected %s but got %s", methods, allowedMethods)
	}

	expectedAllowHeaders := ContentTypeHeader + ", " + AuthorizationHeader
	headers := result.Header.Get("Access-Control-Allow-Headers")
	if headers != expectedAllowHeaders {
		t.Errorf("Wrong Access Control Allow Headers Handler: expected %s but got %s", headers, expectedAllowHeaders)
	}

	exposed := result.Header.Get("Access-Control-Expose-Headers")
	if exposed != AuthorizationHeader {
		t.Errorf("Wrong Access Control Expose Headers Handler: expected %s but got %s", exposed, AuthorizationHeader)
	}

	maxAge := result.Header.Get("Access-Control-Max-Age")
	if maxAge != "600" {
		t.Errorf("Wrong Wrong Access Control Max Age Handler: expected %s but got %s", maxAge, "600")
	}
}
