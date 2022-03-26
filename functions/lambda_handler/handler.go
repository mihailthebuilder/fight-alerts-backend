package handler

import (
	"fight-alerts-backend/scraper"
	"fmt"
)

type Handler struct {
	Scraper scraper.IScraper
}

func (h Handler) HandleRequest() ([]scraper.FightRecord, error) {
	data, err := h.Scraper.GetResultsFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Scraped data:\n%#v\n", data)

	return data, nil
}
