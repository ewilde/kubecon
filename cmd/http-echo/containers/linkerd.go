package containers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/ory-am/dockertest.v3"
)

var linkerdVersion = "1.3.6"

type linkerdContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	uri      string
}

func NewLinkerdContainer(pool *dockertest.Pool, zipkinContainerName string, ipAddress string) (container Container, err error) {

	if err := createServiceDiscoveryFile(ipAddress); err != nil {
		return nil, err
	}

	options := &dockertest.RunOptions{
		Name:         "linkerd",
		Repository:   "buoyantio/linkerd",
		Tag:          linkerdVersion,
		ExposedPorts: []string{"9990", "4140"},
		Links:        []string{zipkinContainerName},
		Mounts:       []string{fmt.Sprintf("%s:/config/", filepath.Join(getCurrentPath(), "containers/linkerd"))},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9990/tcp": {{HostIP: "", HostPort: "9990"}},
			"4140/tcp": {{HostIP: "", HostPort: "4140"}},
		},
		Cmd: []string{"/config/linkerd.config.yml"},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		return nil, err
	}

	linkerdUri := fmt.Sprintf("http://localhost:%v", resource.GetPort("9990/tcp"))
	pool.MaxWait = time.Minute * 1
	if err := pool.Retry(func() error {

		var err error
		if err = checkLinkerdServiceIsStarted(linkerdUri); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to linkerd: %s", err)
	}

	name := getContainerName(resource)
	log.Printf("linkerd %s (%v): up\n", linkerdVersion, name)

	return &linkerdContainer{
		Name:     name,
		pool:     pool,
		resource: resource,
		uri:      linkerdUri,
	}, nil
}
func createServiceDiscoveryFile(ipAddress string) error {
	fileContents := []byte(fmt.Sprintf("%s 5678", ipAddress))
	return ioutil.WriteFile(filepath.Join(getCurrentPath(), "containers/linkerd/disco/service1"), fileContents, 0644)
}

func checkLinkerdServiceIsStarted(linkerdUri string) error {
	client := gorequest.New()
	response, body, err := client.Get(linkerdUri + "/admin/ping").End()

	if err != nil {
		return err[0]
	}

	if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	return nil
}

func getCurrentPath() string {
	path, _ := os.Getwd()
	return path
}

func (l *linkerdContainer) Stop() error {
	return l.pool.Purge(l.resource)
}

func (l *linkerdContainer) GetUri(id string) string {
	return fmt.Sprintf("http://localhost:%s", l.resource.GetPort(id))
}
