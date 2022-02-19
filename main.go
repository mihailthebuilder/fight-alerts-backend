package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Scraper struct {
	url string
}

type IScraper interface {
	getResultsFromUrl() ([]string, error)
}

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func main() {
	fmt.Println("Started")

	var scraper IScraper = Scraper{mmaUrl}

	data, err := scraper.getResultsFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	fmt.Println("Ended")
}

func (s Scraper) getResultsFromUrl() ([]string, error) {
	// Instantiate default collector
	c := colly.NewCollector()

	var results []string
	errOut := fmt.Errorf("unable to find any results")

	c.OnError(func(r *colly.Response, err error) {
		errOut = fmt.Errorf("request URL - %v - failed with status code - %v - error - %v", r.Request.URL, r.StatusCode, err)
	})

	c.OnHTML("#upcoming_tab table tr[onclick]", func(e *colly.HTMLElement) {
		// fmt.Print(e.Text)
		errOut = nil
	})

	c.Visit(s.url)

	return results, errOut
}
