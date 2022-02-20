package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Scraper struct {
	url string
}

type IScraper interface {
	getResultsFromUrl() ([]ProcessedFightRecord, error)
}

func (s Scraper) getResultsFromUrl() ([]ProcessedFightRecord, error) {
	// Instantiate default collector
	c := colly.NewCollector()

	var results []ProcessedFightRecord
	errOut := fmt.Errorf("unable to find any results")

	c.OnError(func(r *colly.Response, err error) {
		errOut = fmt.Errorf("failed with status code - %v - error - %v", r.StatusCode, err)
	})

	c.OnHTML("#upcoming_tab table tr[onclick]", func(e *colly.HTMLElement) {
		var fightRecord IFightRecord = FightRecord{*e}

		processedRecord, err := fightRecord.getProcessedRecord()

		// if we find an error in one of the records, we
		// ignore that record and move on to the next
		if err != nil {
			fmt.Print(err)
			return
		}

		// dateTimeLayout := time.RFC3339
		// dateTimeString := e.ChildAttr("meta[content]", "content")

		// dateTimeParsed, errTimeParse := time.Parse(dateTimeLayout, dateTimeString)

		// if errTimeParse != nil {
		// 	errOut = fmt.Errorf("can't parse string - %v", dateTimeString)
		// 	return
		// }

		// headline := e.ChildText("span[itemprop='name']")

		results = append(results, processedRecord)
		errOut = nil
	})

	c.Visit(s.url)

	return results, errOut
}
