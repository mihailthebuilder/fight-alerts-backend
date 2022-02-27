package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type Scraper struct {
	url string
}

type FightRecord struct {
	DateTime time.Time
	Headline string
}

type IScraper interface {
	getResultsFromUrl() ([]FightRecord, error)
}

func (s Scraper) getResultsFromUrl() ([]FightRecord, error) {
	// Instantiate default collector
	c := colly.NewCollector()

	var results []FightRecord
	errOut := fmt.Errorf("unable to find any results")

	c.OnError(func(r *colly.Response, err error) {
		errOut = fmt.Errorf("failed with status code - %v - error - %v", r.StatusCode, err)
	})

	c.OnHTML("#upcoming_tab table tr[onclick]", func(e *colly.HTMLElement) {
		// html, err := e.DOM.Html()

		// // if we find an error in one of the records, we
		// // ignore that record and move on to the next
		// if err != nil {
		// 	fmt.Print(err)
		// 	return
		// }

		// fmt.Printf("Error processing record: %v\n", html)
		// fmt.Println("Error processing record {unknown}")

		// dateTimeLayout := time.RFC3339
		// dateTimeString := e.ChildAttr("meta[content]", "content")

		// dateTimeParsed, errTimeParse := time.Parse(dateTimeLayout, dateTimeString)

		// if errTimeParse != nil {
		// 	errOut = fmt.Errorf("can't parse string - %v", dateTimeString)
		// 	return
		// }

		// headline := e.ChildText("span[itemprop='name']")

		results = append(results, FightRecord{})
		errOut = nil
	})

	c.Visit(s.url)

	return results, errOut
}
