package main

import (
	"fmt"
)

type handler struct {
	scraper IScraper
}

func (h handler) handleRequest() ([]FightRecord, error) {
	data, err := h.scraper.getResultsFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Scraped data:\n%#v\n", data)

	return data, nil
}
