package main

import "fmt"

type handler struct {
	scraper IScraper
}

func (h handler) handleRequest() {
	fmt.Println("Started")

	data, err := h.scraper.getResultsFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", data)

	fmt.Println("\nEnded")
}
