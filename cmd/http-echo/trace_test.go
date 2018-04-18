package main

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"errors"
	"time"
	"log"
	"encoding/json"
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
	go func() {NewServer("TestTraceViaService", 200, 0, 0)}()
	if err := retry(func() error {
		response, body, err := gorequest.New().
			Get("http://localhost:5678").
			End()

		if err != nil {
			return err[0]
		}

		if response.StatusCode >= 400 {
			return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
		}

		return nil
	}, time.Minute*2); err != nil {
		t.Fatal(err)
	}

	if err := retry(func() error {
		url := fmt.Sprintf("%s/service1", runningContainers["linkerd"].GetUri("4140/tcp"))
		response, body, err := gorequest.New().
			Get(url).
			End()

		if err != nil {
			return err[0]
		}

		if response.StatusCode >= 400 {
			return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
		}

		return nil
	}, time.Minute*2); err != nil {
		t.Fatal(err)
	}

	log.Println("Got response from linkerd")

	var body string
	var err []error
	if _, body, err = gorequest.New().
		Get("http://localhost:9411/api/v2/traces").
		End(); err != nil {
			t.Fatal(err)
	}

	var traces interface{}
	if err := json.Unmarshal([]byte(body), &traces); err != nil {
		t.Fatal(err)
	}

	if data, ok := traces.([]interface{}); ok {
		if firstItem, ok := data[0].([]interface{}); ok {
			if len(firstItem) != 2 {
				t.Errorf("last trace had %d spans expected %d", len(firstItem), 2)
			}

			return
		}
	}

	t.Errorf("Not able to decode trace data")
}
