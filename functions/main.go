package main

import (
	handler "fight-alerts-backend/lambda_handler"
	"fight-alerts-backend/scraper"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	scraper := scraper.Scraper{Url: scraper.MmaUrl}
	lambdaHandler := handler.Handler{Scraper: scraper}
	lambda.Start(lambdaHandler.HandleRequest)
}
