package main

import (
	"bytes"
	"fmt"
	"os"
)

func (s *steps) startContainers() {
	err := s.containers.Start()
	if err != nil {
		panic(err)
	}
}

func (s *steps) stopContainers() {
	fmt.Println("Lambda log:")
	readCloser, err := s.containers.GetLambdaLog()
	if err != nil {
		fmt.Printf("unable to get logs from containers: %v\n", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(readCloser)
	fmt.Println(buf.String())

	fmt.Println("Stopping containers")
	err = s.containers.Stop()
	if err != nil {
		panic(err)
	}
}

func (s *steps) setUpDatastore() {
	fmt.Println("Setting up datastore client...")

	auroraPort, err := s.containers.GetLocalhostPort(s.containers.auroraContainer, AuroraPort)
	if err != nil {
		panic(err)
	}

	hostName := "localhost"
	if os.Getenv("JENKINS") == "true" {
		hostName = "docker"
	}

	s.datastore = AuroraClient{host: hostName, port: auroraPort}

	err = s.datastore.connectToDatabase()
	if err != nil {
		panic(err)
	}

	err = s.datastore.createEventTable()
	if err != nil {
		panic(err)
	}
}

func (s *steps) setUpEventBridgeClient() {
	fmt.Println("Setting up eventbridge client...")

	ebPort, err := s.containers.GetLocalhostPort(s.containers.eventBridgeContainer, MockEventBridgeConxDetails.Port)
	if err != nil {
		panic(err)
	}

	hostName := "localhost"
	if os.Getenv("JENKINS") == "true" {
		hostName = "docker"
	}

	s.eventbridge = EventBridgeClient{host: hostName, port: ebPort}
	err = s.eventbridge.Connect()
	if err != nil {
		panic(err)
	}
}
