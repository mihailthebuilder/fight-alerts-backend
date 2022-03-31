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

func (s *steps) setUpAuroraClient() {
	fmt.Println("Setting up Aurora client")

	auroraPort, err := s.containers.GetLocalhostPort(s.containers.auroraContainer, AuroraPort)
	if err != nil {
		panic(err)
	}
	s.AuroraClient.port = auroraPort
	s.AuroraClient.host = GetHostName()
	s.AuroraClient.dbconx = s.AuroraClient.connectToDatabase()
}

func (s *steps) createAuroraTables() {
	fmt.Println("Creating tables in Aurora")

	q := `
		create table Event (
			EventId int,
			Headline varchar(100),
			DateTime timestamptz
		);
	`

	_, err := s.AuroraClient.dbconx.Exec(q)
	if err != nil {
		panic(err)
	}
}

func GetHostName() string {
	if os.Getenv("JENKINS") == "true" {
		return "docker"
	}
	return "localhost"
}
