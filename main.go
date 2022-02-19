package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func main() {
	fmt.Println("Started")

	data, err := getDataFromUrl(mmaUrl)

	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	fmt.Println("Ended")
}

func getDataFromUrl(url string) ([]string, error) {
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

	c.Visit(url)

	return data, errOut
}
