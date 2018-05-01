package main

import (
	"errors"
	"fmt"
	"github.com/ewilde/kubecon/cmd/http-echo/containers"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
	"testing"
	"time"
)

var runningContainers = map[string]containers.Container{}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	zipkin, err := containers.NewZipkinContainer(pool)
	if err != nil {
		log.Fatalf("Could create zipkin container: %s", err)
	}

	runningContainers["zipkin"] = zipkin

	linkerd, err := containers.NewLinkerdContainer(pool, "zipkin", ipAddress.String())
	if err != nil {
		log.Fatalf("Could create linkerd container: %s", err)
	}

	runningContainers["linkerd"] = linkerd

	code := m.Run()

	linkerd.Stop()
	zipkin.Stop()

	os.Exit(code)
}

func startServer(t *testing.T) {
	go func() { NewServer("TestTraceViaService", 200, 0, 0) }()
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
}
