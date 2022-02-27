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

	dateTime, err := parseDateTime(e.ChildAttr("meta[content]", "content"))

	if err != nil {
		return FightRecord{}, err
	}

	headline := e.ChildText("span[itemprop='name']")
	if len([]rune(headline)) == 0 {
		return FightRecord{}, fmt.Errorf("can't get headline from html")
	}

	return FightRecord{DateTime: dateTime, Headline: headline}, nil
}

func parseDateTime(s string) (time.Time, error) {
	var errOut error = nil

	layout := time.RFC3339

	result, err := time.Parse(layout, s)

	if err != nil {
		errOut = fmt.Errorf("can't parse date from html - %v", err)
	}

	return result, errOut
}
