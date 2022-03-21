package main

import (
	"context"
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
	fmt.Printf("\n\nPORT IS %v\n\n", port)

	return godog.ErrPending
}
