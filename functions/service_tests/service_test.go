package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

var _steps = steps{}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(_steps.startContainers)
	ctx.BeforeSuite(_steps.setUpDatastore)
	ctx.AfterSuite(_steps.stopContainers)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^fight records exist in the database$`, _steps.fightRecordsExistInDb)
	ctx.Step(`^lambda is invoked$`, _steps.lambdaIsInvoked)
	ctx.Step(`^the original fight records are deleted$`, _steps.originalFightRecordsAreDeleted)
	ctx.Step(`^newly-scraped fight records are inserted into the database$`, _steps.newFightRecordsAreInserted)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{"features"},
	}

	status := godog.TestSuite{
		Name:                 "godogs",
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
