package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Scraper struct {
	url string
}

type IScraper interface {
	getResultsFromUrl() ([]FightData, error)
}

type FightData struct {
	Date     string
	Headline string
}

func (s Scraper) getResultsFromUrl() ([]FightData, error) {
	// Instantiate default collector
	c := colly.NewCollector()

	var results []FightData
	errOut := fmt.Errorf("unable to find any results")

	c.OnError(func(r *colly.Response, err error) {
		errOut = fmt.Errorf("request URL - %v - failed with status code - %v - error - %v", r.Request.URL, r.StatusCode, err)
	})

	c.OnHTML("#upcoming_tab table tr[onclick]", func(e *colly.HTMLElement) {
		fightDate := e.ChildAttr("meta[content]", "content")
		headline := e.ChildText("span[itemprop='name']")

		results = append(results, FightData{fightDate, headline})
		errOut = nil
	})

	c.Visit(s.url)

	return results, errOut
}
