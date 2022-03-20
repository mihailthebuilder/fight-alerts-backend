package main

import (
	"context"
	"fmt"
	"io"
	"time"

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
	if err != nil {
		return err
	}

	fmt.Println("Sleeping for 10 seconds while containers start")
	time.Sleep(10 * time.Second)

	return nil
}

func (c *Containers) Stop() error {
	context := context.Background()

	err := c.lambdaContainer.Terminate(context)
	if err != nil {
		return err
	}

	return nil
}

func (c *Containers) startLambdaContainer() error {
	req := testcontainers.ContainerRequest{
		Image:        "lambci/lambda:go1.x",
		ExposedPorts: []string{"9001/tcp"},
		Name:         "lambda",
		Hostname:     "lambda",
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
	err = c.lambdaContainer.CopyFileToContainer(context, "./bin/main", "/var/task/handler", 365)
	if err != nil {
		return err
	}
	c.lambdaContainer.Start(context)

	return nil
}
