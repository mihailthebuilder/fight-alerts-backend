package main

import (
	"fight-alerts-backend/scraper"
	"fmt"
)

type handler struct {
	scraper scraper.IScraper
}

func (h handler) handleRequest() ([]scraper.FightRecord, error) {
	data, err := h.scraper.GetResultsFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Scraped data:\n%#v\n", data)

	return data, nil
}
