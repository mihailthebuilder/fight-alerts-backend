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

type ICollyHtmlElem interface {
	ChildAttr(selector string, attr string) string
	ChildText(selector string) string
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
		record, err := parseCollyHtml(e)

		// if we find an error in one of the records, we
		// ignore that record and move on to the next
		if err != nil {
			// html, err := e.DOM.Html()

			// if err != nil {
			// 	fmt.Print(err)
			// 	html = "unknown"
			// }

			fmt.Print(err)
			return
		}

		results = append(results, record)
	})

	c.Visit(s.url)

	// if len results > 0 return error

	return results, errOut
}

func parseCollyHtml(e ICollyHtmlElem) (FightRecord, error) {

	fightRecord := FightRecord{}
	var errOut error = nil

	dateTimeLayout := time.RFC3339
	dateTimeString := e.ChildAttr("meta[content]", "content")

	dateTimeParsed, errTimeParse := time.Parse(dateTimeLayout, dateTimeString)

	if errTimeParse != nil {
		errOut = fmt.Errorf("can't parse date from html")
	}

	fightRecord.DateTime = dateTimeParsed

	headline := e.ChildText("span[itemprop='name']")
	if len([]rune(headline)) == 0 {
		errOut = fmt.Errorf("can't get headline from html")
	}
	fightRecord.Headline = headline

	return fightRecord, errOut
}
