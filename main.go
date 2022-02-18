package main

import (
	"fmt"
)

const mmaUrl = "https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"

func main() {
	fmt.Println("Started")

	data := getDataFromUrl(mmaUrl)

	fmt.Println(data)

	// // Instantiate default collector
	// c := colly.NewCollector()

	// c.OnError(func(r *colly.Response, err error) {
	// 	fmt.Printf("Request URL: %v\n\nfailed with response: %#v\n\nError: %v\n", r.Request.URL, r, err)
	// })

	// c.Visit("https://abs.fiejfi")

	fmt.Println("Ended")
}

func getDataFromUrl(url string) []string {
	return nil
}
