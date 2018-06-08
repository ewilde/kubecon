package containers

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/ory-am/dockertest.v3"
)

type Container interface {
	Stop() error
	GetUri(id string) string
}

type zipkinContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	uri      string
}

var zipkinVersion = "latest" // "2.6.1"

func NewZipkinContainer(pool *dockertest.Pool) (container Container, err error) {
	envVars := []string{
		"SCRIBE_ENABLED=true",
	}

	options := &dockertest.RunOptions{
		Name:         "zipkin",
		Repository:   "openzipkin/zipkin",
		Tag:          zipkinVersion,
		Env:          envVars,
		ExposedPorts: []string{"9410", "9411"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9410/tcp": {{HostIP: "", HostPort: "9410"}},
			"9411/tcp": {{HostIP: "", HostPort: "9411"}},
		},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	zipkinUri := fmt.Sprintf("http://localhost:%v", resource.GetPort("9411/tcp"))
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
	log.Printf("Zipkin %s (%v): up\n", zipkinVersion, name)

	return &zipkinContainer{
		Name:     name,
		pool:     pool,
		resource: resource,
		uri:      zipkinUri,
	}, nil
}

func getContainerName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}

func checkZipkinServiceIsStarted(zipkinUri string) error {
	client := gorequest.New()
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

func (z *zipkinContainer) GetUri(id string) string {
	return fmt.Sprintf("http://localhost:%s", z.resource.GetPort(id))
}
