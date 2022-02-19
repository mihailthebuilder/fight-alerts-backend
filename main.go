package main

import (
	"fmt"
)

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func main() {
	fmt.Println("Started")

	var scraper IScraper = Scraper{mmaUrl}

	data, err := scraper.getResultsFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	fmt.Println("Ended")
}
