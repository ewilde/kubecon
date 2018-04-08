package main

import (
	"log"
	"fmt"
	"time"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/ory-am/dockertest.v3"
	"strings"
	"errors"
)

type container interface {
	Stop() error
}

type zipkinContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	Uri      string
}

func newZipkinContainer(pool *dockertest.Pool, ) (container *container, err error) {
	envVars := []string{
		"SCRIBE=true",
	}

	options := &dockertest.RunOptions{
		Name:         "zipkin",
		Repository:   "openzipkin/zipkin:2.6.1",
		Env:          envVars,
		ExposedPorts: []string{"5601"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5601/tcp": {{HostIP: "", HostPort: "5601"}},
		},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	zipkinUri := fmt.Sprintf("http://localhost:%v", resource.GetPort("5601/tcp"))
	pool.MaxWait = time.Minute * 1
	if err := pool.Retry(func() error {

		var err error
		if err = checkZipkinServiceIsStarted(zipkinUri); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to kibana: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not connect to kibana: %s", err)
	}

	name := getContainerName(resource)
	log.Printf("Kibana %s (%v): up\n", kibanaVersion, name)

	return &zipkinContainer{
		Name:     name,
		pool:     pool,
		resource: resource,
		Uri:      zipkinUri,
	}, nil
}

func getContainerName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}

func checkZipkinServiceIsStarted(client **gorequest.SuperAgent, zipkinUri string) error {
	response, body, err := client.Get(zipkinUri + "/health").End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func (z *zipkinContainer) Stop() error {
	return z.pool.Purge(z.resource)
}
