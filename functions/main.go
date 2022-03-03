package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func main() {
	scraper := Scraper{mmaUrl}
	lambdaHandler := handler{scraper: scraper}
	lambda.Start(lambdaHandler.handleRequest)
}

func run(s IScraper) {
	fmt.Println("Started")

	data, err := s.getResultsFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", data)

	fmt.Println("\nEnded")
}
