package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Scraper struct {
	url string
}

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func main() {
	fmt.Println("Started")

	scraper := Scraper{mmaUrl}

	data, err := scraper.getDataFromUrl()

	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	fmt.Println("Ended")
}

func (s *Scraper) getDataFromUrl() ([]string, error) {
	// Instantiate default collector
	c := colly.NewCollector()

	var data []string
	errOut := fmt.Errorf("something unknown went wrong")

	c.OnError(func(r *colly.Response, err error) {
		errOut = fmt.Errorf("request URL - %v - failed with status code - %v - error - %v", r.Request.URL, r.StatusCode, err)
	})

	c.OnHTML("#upcoming_tab table tr[onclick]", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		errOut = nil
	})

	c.Visit(s.url)

	return data, errOut
}
