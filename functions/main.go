package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func main() {
	scraper := Scraper{mmaUrl}
	lambdaHandler := handler{scraper: scraper}
	lambda.Start(lambdaHandler.handleRequest)
}
