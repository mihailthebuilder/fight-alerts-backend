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
	RawRecord colly.HTMLElement
}

type IFightRecord interface {
	getHtmlString() string
	getProcessedRecord() (ProcessedFightRecord, error)
}

func (f FightRecord) getProcessedRecord() (ProcessedFightRecord, error) {
	var processedFightRecord ProcessedFightRecord

	return processedFightRecord, nil
}

func (f FightRecord) getHtmlString() string {
	html, err := f.RawRecord.DOM.Html()
	if err != nil {
		return "{unknown}"
		// fmt.Println("Error processing record {unknown}")
	}

	// fmt.Printf("Error processing record: %v\n", html)
	return html
}
