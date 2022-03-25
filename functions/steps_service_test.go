package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type steps struct {
	containers Containers
}

func (s *steps) sherdogIsAvailable() error {
	resp, err := http.Get(mmaUrl)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("request didn't work - %#v", resp)
	}

	return nil
}

type lambdaPortKey string

func (s *steps) lambdaIsInvoked(ctx context.Context) (context.Context, error) {

	port, err := s.containers.GetLocalHostLambdaPort()

	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, lambdaPortKey("lambdaPort"), port)
	return ctx, nil
}

func (s *steps) fightDataIsReturned(ctx context.Context) error {
	port := ctx.Value(lambdaPortKey("lambdaPort"))

	host := "localhost"
	if os.Getenv("JENKINS") == "true" {
		host = "docker"
	}

	url := fmt.Sprintf("http://%s:%d/2015-03-31/functions/myfunction/invocations", host, port)

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

	var records []FightRecord
	json.Unmarshal([]byte(body), &records)

	if len(records) < 3 {
		return fmt.Errorf("can't scrape fight records %#v", records)
	}

	return nil
}
