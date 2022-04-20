package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type steps struct {
	containers Containers
	datastore  AuroraClient
}

func (s *steps) sherdogIsAvailable() error {
	mmaUrl := "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

	resp, err := http.Get(mmaUrl)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("request didn't work - %#v", resp)
	}

	return nil
}

func (s *steps) lambdaIsInvoked(ctx context.Context) error {

	port, err := s.containers.GetLocalhostPort(s.containers.lambdaContainer, LambdaPort)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d/2015-03-31/functions/myfunction/invocations", GetHostName(), port)

	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body := buf.String()

	if response.StatusCode != 200 {
		return fmt.Errorf("invoking Lambda: %d %s", response.StatusCode, body)
	}

	return nil
}

func (s *steps) scrapedDataIsInDb(ctx context.Context) error {

	rows, err := s.datastore.getAllRowsFromEventTable()

	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return fmt.Errorf("no items in datastore: %v", rows)
	}

	return nil
}
