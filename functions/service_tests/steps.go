package main

import (
	"net/http"

	"github.com/cucumber/godog"
)

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func sherdogIsAvailableForAccess() error {
	resp, err := http.Get(mmaUrl)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}

	return nil
}

func theServiceIsInvoked() error {
	return godog.ErrPending
}

func allTheFightDataIsLogged() error {
	return godog.ErrPending
}
