package main

import (
	"time"

	"github.com/gocolly/colly"
)

type ProcessedFightRecord struct {
	DateTime time.Time
	Headline string
}

type FightRecord struct {
	RawRecord *colly.HTMLElement
	ProcessedFightRecord
}

func (f *FightRecord) process() error {
	return nil

	// return fmt.Errorf("unable to process")
}

func (f *FightRecord) convertString() string {
	html, err := f.RawRecord.DOM.Html()
	if err != nil {
		return "{unknown}"
		// fmt.Println("Error processing record {unknown}")
	}

	// fmt.Printf("Error processing record: %v\n", html)
	return html
}
