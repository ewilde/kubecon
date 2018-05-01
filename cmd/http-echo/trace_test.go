package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"testing"
	"time"
)

func TestTraceViaLinkerd(t *testing.T) {
	startServer(t)
	callServiceViaLinkerd(t)
	foundTrace := getSpanFromZipkin(t)


	if foundTrace == nil {
		t.Errorf("Could not find a trace with the expected 3 spans")
	}

	span := foundTrace[0].(map[string]interface{})
	localEndpoint := span["localEndpoint"].(map[string]interface{})
	serviceName := localEndpoint["serviceName"]

	if serviceName != "http-echo" {
		t.Errorf("Expected service name to be http-echo, actual %s", serviceName)
	}
}

func getSpanFromZipkin(t *testing.T) []interface{} {
	var foundTrace []interface{}
	if err := retry(func() error {
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
			for _, item := range data {
				dataItem := item.([]interface{})
				if len(dataItem) == 3 {
					foundTrace = dataItem
				}
			}
		}

		if foundTrace == nil {
			return errors.New("trace with 3 spans not found yet")
		}

		return nil
	}, time.Minute*2); err != nil {
		t.Error(err)
	}

	return foundTrace
}

func callServiceViaLinkerd(t *testing.T) {
	if err := retry(func() error {
		url := fmt.Sprintf("%s/service1", runningContainers["linkerd"].GetUri("4140/tcp"))
		response, body, err := gorequest.New().
			Get(url).
			End()

		if err != nil {
			log.Println(err)
			return err[0]
		}

		if response.StatusCode >= 400 {
			err := errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
			log.Println(err)
			return err
		}

		return nil
	}, time.Minute*2); err != nil {
		t.Fatal(err)
	}
}
