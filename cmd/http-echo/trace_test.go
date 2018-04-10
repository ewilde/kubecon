package main

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTraceViaService(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(httpLog(stdoutW, httpEcho(stdoutW, "serviceA", 200, 0, 0)))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `serviceA`
	want := strings.TrimSpace(rr.Body.String())
	if want != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", want, expected)

	}
}

func TestTraceViaLinkerd(t *testing.T) {
	client := gorequest.New()
	response, body, err := client.Get("http://localhost:4140/service1").End()

	if err != nil {
		t.Fatal(err[0])
	}

	if response.StatusCode >= 400 {
		t.Fatal(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}
}
