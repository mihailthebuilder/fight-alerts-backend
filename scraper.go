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
	getResultsFromUrl() ([]ProcessedFightRecord, error)
}

type ProcessedFightRecord struct {
	DateTime time.Time
	Headline string
}

// type FightRecord struct {
// 	RawRecord *colly.HTMLElement
// 	ProcessedFightRecord
// }

func (s Scraper) getResultsFromUrl() ([]ProcessedFightRecord, error) {
	// Instantiate default collector
	c := colly.NewCollector()

	var results []ProcessedFightRecord
	errOut := fmt.Errorf("unable to find any results")

	c.OnError(func(r *colly.Response, err error) {
		errOut = fmt.Errorf("failed with status code - %v - error - %v", r.StatusCode, err)
	})

	c.OnHTML("#upcoming_tab table tr[onclick]", func(e *colly.HTMLElement) {
		dateTimeLayout := time.RFC3339
		dateTimeString := e.ChildAttr("meta[content]", "content")

		dateTimeParsed, errTimeParse := time.Parse(dateTimeLayout, dateTimeString)

		if errTimeParse != nil {
			errOut = fmt.Errorf("can't parse string - %v", dateTimeString)
			return
		}

		headline := e.ChildText("span[itemprop='name']")

		results = append(results, ProcessedFightRecord{dateTimeParsed, headline})
		errOut = nil
	})

	c.Visit(s.url)

	return results, errOut
}
