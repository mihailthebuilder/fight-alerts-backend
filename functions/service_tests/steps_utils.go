package main

import (
	"bytes"
	"fight-alerts-backend/datastore"
	utils "fight-alerts-backend/test_utils"
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
	fmt.Println("Setting up datastore")

	auroraPort, err := s.containers.GetLocalhostPort(s.containers.auroraContainer, AuroraPort)
	if err != nil {
		panic(err)
	}

	s.datastore = &datastore.Datastore{
		Host:     GetHostName(),
		Port:     auroraPort,
		User:     PostgresConxDetails.User,
		Password: PostgresConxDetails.Password,
		Dbname:   PostgresConxDetails.Database,
	}

	err = s.datastore.Connect()
	if err != nil {
		panic(err)
	}

	err = utils.CreateEventTable(s.datastore.Db)
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
