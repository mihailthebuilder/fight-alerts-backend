package main

import (
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

type steps struct {
	containers Containers
}

func (s *steps) sherdogIsAvailable() error {
	resp, err := http.Get(mmaUrl)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}

	return nil
}

func (s *steps) lambdaIsInvoked() error {

	localLambdaInvocationPort, err := s.containers.GetLocalHostLambdaPort()
	if err != nil {
		return err
	}

	fmt.Printf("\nLambda active on port %v\n", localLambdaInvocationPort)

	return nil
}

func (s *steps) fightDataIsLogged() error {
	return godog.ErrPending
}
