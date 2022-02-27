package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockScraper struct {
	wantError bool
}

func (s MockScraper) getResultsFromUrl() ([]FightRecord, error) {
	if s.wantError == true {
		return nil, fmt.Errorf("Panic example")
	}
	return []FightRecord{}, nil
}

func Test_run(t *testing.T) {
	type args struct {
		s IScraper
	}

	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{name: "Test run function successful", args: args{s: MockScraper{wantError: false}}},
		{name: "Test run function failed", args: args{s: MockScraper{wantError: true}}, wantPanic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic == true {
				assert.Panics(t, func() { run(tt.args.s) }, "The code did not panic")
				return
			}

			run(tt.args.s)
		})
	}
}
