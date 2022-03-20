package main

import (
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
	_, err := s.containers.GetLocalHostLambdaPort()
	return err
}

func (s *steps) fightDataIsLogged() error {
	return godog.ErrPending
}
