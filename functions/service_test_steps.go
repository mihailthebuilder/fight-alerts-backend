package main

import (
	"net/http"

	"github.com/cucumber/godog"
)

func sherdogIsAvailable() error {
	resp, err := http.Get(mmaUrl)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}

	return nil
}

func lambdaIsInvoked() error {
	return godog.ErrPending
}

func fightDataIsLogged() error {
	return godog.ErrPending
}
