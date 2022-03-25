package main

import (
	"fight-alerts-backend/resources"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	scraper := Scraper{resources.MmaUrl}
	lambdaHandler := handler{scraper: scraper}
	lambda.Start(lambdaHandler.handleRequest)
}
