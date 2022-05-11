package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

func portStringBuilder(port int) string {
	return fmt.Sprintf("%d/tcp", port)
}

var LambdaPort = 9001
var AuroraPort = 5432
var MyNetwork = "myNetwork"

type DbConxDetails struct {
	Password string
	User     string
	Database string
	Host     string
}

var PostgresConxDetails = DbConxDetails{
	Password: "password",
	User:     "FightAlertsUser",
	Database: "FightAlertsDb",
	Host:     "postgres",
}

type EventBridgeConxDetails struct {
	Port int
	Host string
}

var MockEventBridgeConxDetails = EventBridgeConxDetails{
	Port: 4566,
	Host: "eventbridge",
}

type Containers struct {
	network              testcontainers.Network
	lambdaContainer      testcontainers.Container
	auroraContainer      testcontainers.Container
	eventBridgeContainer testcontainers.Container
}

func (c *Containers) GetLambdaLog() (io.ReadCloser, error) {
	return c.lambdaContainer.Logs(context.Background())
}

func (c *Containers) Start() error {
	var err error

	err = c.createNetwork()
	if err != nil {
		return err
	}

	err = c.startLambdaContainer()
	if err != nil {
		return err
	}

	err = c.startAuroraContainer()
	if err != nil {
		return err
	}

	err = c.startEventBridgeContainer()
	if err != nil {
		return err
	}

	fmt.Println("Sleeping for 5 seconds while containers start")
	time.Sleep(5 * time.Second)

	return nil
}

func (c *Containers) Stop() error {
	context := context.Background()

	var err error

	err = c.lambdaContainer.Terminate(context)
	if err != nil {
		return err
	}

	err = c.auroraContainer.Terminate(context)
	if err != nil {
		return err
	}

	err = c.eventBridgeContainer.Terminate(context)
	if err != nil {
		return err
	}

	err = c.network.Remove(context)
	if err != nil {
		return err
	}

	return nil
}

func (c *Containers) startLambdaContainer() error {
	req := testcontainers.ContainerRequest{
		Image:        "lambci/lambda:go1.x",
		ExposedPorts: []string{portStringBuilder(LambdaPort)},
		Name:         "lambda",
		Hostname:     "lambda",
		Env: map[string]string{
			"DOCKER_LAMBDA_STAY_OPEN":  "1",
			"AWS_ACCESS_KEY_ID":        "x",
			"AWS_SECRET_ACCESS_KEY":    "x",
			"RDS_HOST":                 PostgresConxDetails.Host,
			"RDS_USERNAME":             PostgresConxDetails.User,
			"RDS_PASSWORD":             PostgresConxDetails.Password,
			"NOTIFICATION_LAMBDA_ARN":  "arn:aws:lambda:us-east-1:111111111111:function:mock-lambda-arn",
			"EVENTS_ENDPOINT_OVERRIDE": fmt.Sprintf("http://%s:%d", MockEventBridgeConxDetails.Host, MockEventBridgeConxDetails.Port),
		},
		Networks:    []string{MyNetwork},
		NetworkMode: container.NetworkMode(MyNetwork),
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
	err = c.lambdaContainer.CopyFileToContainer(context, "scraper-lambda", "/var/task/handler", 365)
	if err != nil {
		return err
	}
	c.lambdaContainer.Start(context)

	return nil
}

func (c *Containers) GetLocalhostPort(container testcontainers.Container, port int) (int, error) {
	context := context.Background()
	mappedPort, err := container.MappedPort(context, nat.Port(portStringBuilder(port)))
	if err != nil {
		return 0, err
	}
	return mappedPort.Int(), nil
}

func (c *Containers) startAuroraContainer() error {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{portStringBuilder(AuroraPort)},
		Name:         "postgres",
		Hostname:     PostgresConxDetails.Host,
		Env: map[string]string{
			"POSTGRES_PASSWORD": PostgresConxDetails.Password,
			"POSTGRES_USER":     PostgresConxDetails.User,
			"POSTGRES_DB":       PostgresConxDetails.Database,
		},
		Networks:    []string{MyNetwork},
		NetworkMode: container.NetworkMode(MyNetwork),
	}

	var err error

	context := context.Background()

	c.auroraContainer, err = testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		return err
	}

	err = c.auroraContainer.Start(context)
	if err != nil {
		return err
	}

	return nil
}

func (c *Containers) startEventBridgeContainer() error {
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack",
		ExposedPorts: []string{portStringBuilder(MockEventBridgeConxDetails.Port)},
		Name:         "eventbridge",
		Hostname:     MockEventBridgeConxDetails.Host,
		Env: map[string]string{
			"SERVICES":       "events",
			"DEFAULT_REGION": "us-east-1",
			"DEBUG":          "1",
		},
		Networks:    []string{MyNetwork},
		NetworkMode: container.NetworkMode(MyNetwork),
	}

	var err error
	c.eventBridgeContainer, err = testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Containers) createNetwork() error {
	context := context.Background()
	var err error
	req := testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Driver: "bridge", Name: MyNetwork, Attachable: true},
	}
	c.network, err = testcontainers.GenericNetwork(context, req)
	if err != nil {
		return err
	}
	return nil
}
