package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
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
	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", port)

	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		body := buf.String()
		return fmt.Errorf("invoking Lambda: %d %s", response.StatusCode, body)
	}

	return nil
}
