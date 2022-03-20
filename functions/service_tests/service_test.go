package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{"features"},
	}

	status := godog.TestSuite{
		Name: "godogs",
		// TestSuiteInitializer: InitializeTestSuite,
		// ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
