package main

import (
	"bytes"
	"context"
	"fight-alerts-backend/scraper"
	"fmt"
	"net/http"
)

type steps struct {
	containers   Containers
	AuroraClient AuroraClient
}

func (s *steps) sherdogIsAvailable() error {
	resp, err := http.Get(scraper.MmaUrl)

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
	items := s.AuroraClient.getAllItems()
	if len(items) == 0 {
		return fmt.Errorf("no items in datastore: %v", items)
	}

	return nil
}
