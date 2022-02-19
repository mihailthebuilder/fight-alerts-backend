package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type Scraper struct {
	url string
}

type IScraper interface {
	getResultsFromUrl() ([]FightData, error)
}

type FightData struct {
	DateTime time.Time
	Headline string
}

func (s Scraper) getResultsFromUrl() ([]FightData, error) {
	// Instantiate default collector
	c := colly.NewCollector()

	var results []FightData
	errOut := fmt.Errorf("unable to find any results")

	c.OnError(func(r *colly.Response, err error) {
		errOut = fmt.Errorf("failed with status code - %v - error - %v", r.StatusCode, err)
	})

	c.OnHTML("#upcoming_tab table tr[onclick]", func(e *colly.HTMLElement) {
		dateTimeLayout := time.RFC822
		dateTimeString := e.ChildAttr("meta[content]", "content")

		dateTimeParsed, errTime := time.Parse(dateTimeLayout, dateTimeString)

		if errTime != nil {
			errOut = fmt.Errorf("can't parse string %v", dateTimeString)
			return
		}

		headline := e.ChildText("span[itemprop='name']")

		results = append(results, FightData{dateTimeParsed, headline})
		errOut = nil
	})

	c.Visit(s.url)

	return results, errOut
}
