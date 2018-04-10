package containers

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/ory-am/dockertest.v3"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

var linkerdVersion = "1.3.6"

type linkerdContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	Uri      string
}

func NewLinkerdContainer(pool *dockertest.Pool, zipkinContainerName string) (container container, err error) {

	if err := createServiceDiscoveryFile(); err != nil {
		return nil, err
	}

	options := &dockertest.RunOptions{
		Name:         "linkerd",
		Repository:   "buoyantio/linkerd",
		Tag:          linkerdVersion,
		ExposedPorts: []string{"9990,4140"},
		Links:        []string{zipkinContainerName},
		Mounts:       []string{fmt.Sprintf("%s:/config/", filepath.Join(getCurrentPath(), "cmd/http-echo/containers/linkerd"))},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9990/tcp": {{HostIP: "", HostPort: "9990"}},
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
		log.Fatalf("Could not connect to kibana: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not connect to kibana: %s", err)
	}

	name := getContainerName(resource)
	log.Printf("linkerd %s (%v): up\n", linkerdVersion, name)

	return &linkerdContainer{
		Name:     name,
		pool:     pool,
		resource: resource,
		Uri:      linkerdUri,
	}, nil
}
func createServiceDiscoveryFile() error {
	fileContents := []byte(fmt.Sprintf("%s 80", getOutboundIP().String()))
	return ioutil.WriteFile(filepath.Join(getCurrentPath(), "deployments/local/system-2/linkerd/disco/service1"), fileContents, 0644)
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("[WARN] error closing connection for ip check %v", err)
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func (z *linkerdContainer) Stop() error {
	return z.pool.Purge(z.resource)
}
