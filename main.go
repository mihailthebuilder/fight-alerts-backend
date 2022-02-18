package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("Started")

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %v\n\nfailed with response: %#v\n\nError: %v\n", r.Request.URL, r, err)
	})

	c.Visit("https://abs.fiejfi")

	fmt.Println("Ended")
}
