package main

import (
	"fight-alerts-backend/resources"
	"fight-alerts-backend/scraper"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	scraper := scraper.Scraper{resources.MmaUrl}
	lambdaHandler := handler{scraper: scraper}
	lambda.Start(lambdaHandler.handleRequest)
}
