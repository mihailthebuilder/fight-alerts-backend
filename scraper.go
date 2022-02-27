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
	ChildAttr(string, string) string
	ChildText(string) string
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
		html, err := e.DOM.Html()

		// if we find an error in one of the records, we
		// ignore that record and move on to the next
		if err != nil {
			fmt.Print(err)
			html = "unknown"
		}

		parseCollyHtml(e, &results, &errOut, html)
	})

	c.Visit(s.url)

	return results, errOut
}

func parseCollyHtml(e ICollyHtmlElem, r *[]FightRecord, err *error, rawHtml string) {
	dateTimeLayout := time.RFC3339
	dateTimeString := e.ChildAttr("meta[content]", "content")

	dateTimeParsed, errTimeParse := time.Parse(dateTimeLayout, dateTimeString)

	if errTimeParse != nil {
		fmt.Printf("can't parse date from html - %v", rawHtml)
		return
	}

	headline := e.ChildText("span[itemprop='name']")

	*r = append(*r, FightRecord{DateTime: dateTimeParsed, Headline: headline})
	*err = nil
}
