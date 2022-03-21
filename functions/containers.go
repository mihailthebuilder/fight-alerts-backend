package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

type Containers struct {
	lambdaContainer testcontainers.Container
}

func (c *Containers) GetLambdaLog() (io.ReadCloser, error) {
	return c.lambdaContainer.Logs(context.Background())
}

func (c *Containers) Start() error {
	err := c.startLambdaContainer()

	fmt.Println("Sleeping for 5 seconds while containers start")
	time.Sleep(5 * time.Second)

	return err
}

func (c *Containers) Stop() error {
	context := context.Background()

	err := c.lambdaContainer.Terminate(context)
	return err
}

func (c *Containers) startLambdaContainer() error {
	req := testcontainers.ContainerRequest{
		Image:        "lambci/lambda:go1.x",
		ExposedPorts: []string{"9001/tcp"},
		Name:         "lambda",
		Hostname:     "lambda",
		Env: map[string]string{
			"DOCKER_LAMBDA_STAY_OPEN": "1",
			"AWS_ACCESS_KEY_ID":       "x",
			"AWS_SECRET_ACCESS_KEY":   "x",
		},
	}
	context := context.Background()

	var err error
	c.lambdaContainer, err = testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		return err
	}
	err = c.lambdaContainer.CopyFileToContainer(context, "bin/scraper", "/var/task/handler", 365)
	if err != nil {
		return err
	}
	c.lambdaContainer.Start(context)

	return nil
}

func (c *Containers) GetLocalHostLambdaPort() (int, error) {
	return c.GetLocalhostPort(c.lambdaContainer, 9001)
}

func (c *Containers) GetLocalhostPort(container testcontainers.Container, port int) (int, error) {
	context := context.Background()
	mappedPort, err := container.MappedPort(context, nat.Port(fmt.Sprintf("%d/tcp", port)))
	if err != nil {
		return 0, err
	}
	return mappedPort.Int(), nil
}
