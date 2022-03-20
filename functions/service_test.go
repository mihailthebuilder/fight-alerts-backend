package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Sherdog is available for access$`, sherdogIsAvailableForAccess)
	ctx.Step(`^the service is invoked$`, theServiceIsInvoked)
	ctx.Step(`^all the fight data is logged$`, allTheFightDataIsLogged)
}

var _steps = steps{}

type steps struct {
	containers Containers
}

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

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(_steps.startContainers)
	ctx.AfterSuite(_steps.stopContainers)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{"features"},
	}

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
