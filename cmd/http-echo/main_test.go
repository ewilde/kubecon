package main

import (
	//"github.com/ewilde/kubecon/cmd/http-echo/containers"
	//"gopkg.in/ory-am/dockertest.v3"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestRandom_float(t *testing.T) {
	rate := 0.1
	delay := 1000.0

	for i := 0; i < 10000; i++ {
		if rand.Float64() <= rate/100 {
			duration := time.Duration(float64(time.Millisecond) * delay)
			log.Printf("[INFO] Will delay for %s.", duration.String())
		}
	}
}

func TestMain(m *testing.M) {
	//pool, err := dockertest.NewPool("")
	//if err != nil {
	//	log.Fatalf("Could not connect to docker: %s", err)
	//}
	//
	//zipkin, err := containers.NewZipkinContainer(pool)
	//if err != nil {
	//	log.Fatalf("Could create zipkin container: %s", err)
	//}
	//
	//linkerd, err := containers.NewLinkerdContainer(pool, "zipkin")
	//if err != nil {
	//	log.Fatalf("Could create linkerd container: %s", err)
	//}

	code := m.Run()

	//linkerd.Stop()
	//zipkin.Stop()

	os.Exit(code)
}
