package main

import (
	"fight-alerts-backend/scraper"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockScraper struct {
	wantError bool
}

func (s MockScraper) GetResultsFromUrl() ([]scraper.FightRecord, error) {
	if s.wantError == true {
		return nil, fmt.Errorf("Panic example")
	}
	return []scraper.FightRecord{}, nil
}

func Test_handleRequest(t *testing.T) {
	type args struct {
		s scraper.IScraper
	}

	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{name: "run function should not return panic", args: args{s: MockScraper{wantError: false}}},
		{name: "run function should return panic", args: args{s: MockScraper{wantError: true}}, wantPanic: true},
	}
	for _, tt := range tests {
		lambdaHandler := handler{scraper: tt.args.s}

		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic == true {
				assert.Panics(t, func() { lambdaHandler.handleRequest() }, "The code did not panic")
				return
			}

			lambdaHandler.handleRequest()
		})
	}
}
