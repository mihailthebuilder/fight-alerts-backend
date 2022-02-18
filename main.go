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
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", *r, "\nError:", err)
	})

	c.Visit("https://abs.fiejfi")

	fmt.Println("Ended")
}
