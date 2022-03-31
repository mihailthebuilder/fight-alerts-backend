package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

var _steps = steps{}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(_steps.startContainers)
	ctx.BeforeSuite(_steps.setUpAuroraClient)
	ctx.BeforeSuite(_steps.createAuroraTables)
	ctx.AfterSuite(_steps.stopContainers)
	ctx.AfterSuite(_steps.AuroraClient.closeDatabaseConnexion)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^Sherdog is available for access$`, _steps.sherdogIsAvailable)
	ctx.Step(`^lambda is invoked$`, _steps.lambdaIsInvoked)
	ctx.Step(`^scraped fight data is available in the database$`, _steps.scrapedDataIsInDb)
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
