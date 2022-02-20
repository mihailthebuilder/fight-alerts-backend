package main

import (
	"fmt"
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
	return fmt.Errorf("unable to process")
}
