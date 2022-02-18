package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("Starting")

	// Instantiate default collector
	c := colly.NewCollector()
	c.Visit("https://abs")
}
